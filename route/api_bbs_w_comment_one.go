package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_comment_one(db *sql.DB, config tool.Config, already_auth_check bool) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    sub_code := other_set["sub_code"]
    sub_code_parts := strings.Split(sub_code, "-")
    sub_code_last := ""
    new_sub_code := ""

    if other_set["tool"] == "around" {
        new_sub_code = other_set["sub_code"]
    } else {
        if len(sub_code_parts) > 2 {
            sub_code_last = sub_code_parts[len(sub_code_parts)-1]
            sub_code_parts = sub_code_parts[:len(sub_code_parts)-1]

            new_sub_code = strings.Join(sub_code_parts, "-")
        }
    }

    var rows *sql.Rows
    if other_set["tool"] == "around" {
        stmt, err := db.Prepare(tool.DB_change("select set_name, set_data, set_code, set_id from bbs_data where (set_name = 'comment' or set_name like 'comment%') and set_id = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        rows, err = stmt.Query(new_sub_code)
        if err != nil {
            panic(err)
        }
    } else {
        stmt, err := db.Prepare(tool.DB_change("select set_name, set_data, set_code, set_id from bbs_data where (set_name = 'comment' or set_name like 'comment%') and set_id = ? and set_code = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        rows, err = stmt.Query(new_sub_code, sub_code_last)
        if err != nil {
            panic(err)
        }
    }
    defer rows.Close()

    data_list := []map[string]string{}
    temp_dict := map[string]string{}
    ip_parser_temp := map[string][]string{}
    before_set_code := ""

    for rows.Next() {
        var set_name string
        var set_data string
        var set_code string
        var set_id string

        err := rows.Scan(&set_name, &set_data, &set_code, &set_id)
        if err != nil {
            panic(err)
        }

        if before_set_code != set_code {
            if before_set_code != "" {
                data_list = append(data_list, temp_dict)
            }

            temp_dict = map[string]string{}
            temp_dict["id"] = set_id
            temp_dict["code"] = set_code

            before_set_code = set_code
        }

        if set_name == "comment_user_id" {
            var ip_pre string
            var ip_render string

            if _, ok := ip_parser_temp[set_data]; ok {
                ip_pre = ip_parser_temp[set_data][0]
                ip_render = ip_parser_temp[set_data][1]
            } else {
                ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
                ip_render = tool.IP_parser(db, set_data, config.IP)

                ip_parser_temp[set_data] = []string{ip_pre, ip_render}
            }

            temp_dict["comment_user_id"] = ip_pre
            temp_dict["comment_user_id_render"] = ip_render
        } else {
            temp_dict[set_name] = set_data
        }
    }

    if before_set_code != "" {
        data_list = append(data_list, temp_dict)
    }

    return_data := make(map[string]interface{})
    if !already_auth_check {
        if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
            data_list = []map[string]string{}
            return_data["response"] = "require auth"
        }
    }

    if other_set["legacy"] != "" {
        var json_data []byte
        if other_set["tool"] == "around" {
            json_data, _ = json.Marshal(data_list)
        } else {
            json_data, _ = json.Marshal(data_list[0])
        }

        return string(json_data)
    } else {
        return_data["data"] = data_list

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
}
