package route

import (
	"database/sql"
	"encoding/json"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_comment_all(db *sql.DB, sub_code string, already_auth_check bool) []map[string]string {
    end_data := []map[string]string{}

    inter_other_set := map[string]string{}
    inter_other_set["sub_code"] = sub_code
    inter_other_set["tool"] = "around"
    inter_other_set["legacy"] = "on"

    json_data, _ := json.Marshal(inter_other_set)

    send_request := tool.Config{
        Other_set: string(json_data),
        IP: "",
        Cookies: "",
    }
    
    return_data := Api_bbs_w_comment_one(db, send_request, already_auth_check)

    return_data_api := []map[string]string{}
    json.Unmarshal([]byte(return_data), &return_data_api)

    for for_a := 0; for_a < len(return_data_api); for_a++ {
        end_data = append(end_data, return_data_api[for_a])

        temp := Api_bbs_w_comment_all(db, sub_code + "-" + return_data_api[for_a]["code"], already_auth_check)
        if len(temp) > 0 {
            for for_b := 0; for_b < len(temp); for_b++ {
                end_data = append(end_data, temp[for_b])
            }
        }
    }

    return end_data
}

func Api_bbs_w_comment(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    if other_set["tool"] == "length" {
        stmt, err := db.Prepare(tool.DB_change("select count(*) from bbs_data where set_name = 'comment_date' and set_id = ? order by set_code + 0 desc"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        var comment_length string
        bbs_and_post_num := other_set["sub_code"]

        err = stmt.QueryRow(bbs_and_post_num).Scan(&comment_length)
        if err != nil {
            if err == sql.ErrNoRows {
                comment_length = "0"
            } else {
                panic(err)
            }
        }

        stmt, err = db.Prepare(tool.DB_change("select count(*) from bbs_data where set_name = 'comment_date' and set_id like ? order by set_code + 0 desc"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        var reply_length string

        err = stmt.QueryRow(bbs_and_post_num + "-%").Scan(&reply_length)
        if err != nil {
            if err == sql.ErrNoRows {
                reply_length = "0"
            } else {
                panic(err)
            }
        }

        comment_length_int, _ := strconv.Atoi(comment_length)
        reply_length_int, _ := strconv.Atoi(reply_length)

        length_int := comment_length_int + reply_length_int
        length_str := strconv.Itoa(length_int)

        data_list := map[string]string{
            "comment": comment_length,
            "reply":   reply_length,
            "data":    length_str,
        }

        json_data, _ := json.Marshal(data_list)
        return string(json_data)
    } else {
        return_data := make(map[string]interface{})
        
        temp := []map[string]string{}
        if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
            return_data["response"] = "require auth"
        } else {
            temp = Api_bbs_w_comment_all(db, other_set["sub_code"], true)
        }

        if other_set["legacy"] != "" {
            json_data, _ := json.Marshal(temp)
            return string(json_data)
        } else {
            return_data["language"] = map[string]string{
                "normal" : tool.Get_language(db, "normal", false),
                "comment" : tool.Get_language(db, "comment", false),
                "tool" : tool.Get_language(db, "tool", false),
                "return" : tool.Get_language(db, "return", false),
                "upvote" : tool.Get_language(db, "upvote", false),    
            }
            return_data["data"] = temp

            json_data, _ := json.Marshal(return_data)
            return string(json_data)
        }
    }
}
