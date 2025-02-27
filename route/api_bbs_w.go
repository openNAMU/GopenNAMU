package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    stmt, err := db.Prepare(tool.DB_change("select set_name, set_data from bbs_data where set_id = ? and set_code = ?"))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    sub_code := other_set["sub_code"]
    sub_code_parts := strings.Split(sub_code, "-")

    bbs_num := ""
    post_num := ""

    if len(sub_code_parts) > 1 {
        bbs_num = sub_code_parts[0]
        post_num = sub_code_parts[1]
    }

    rows, err := stmt.Query(bbs_num, post_num)
    if err != nil {
        panic(err)
    }

    data_list := map[string]string{}

    for rows.Next() {
        var set_name string
        var set_data string

        err := rows.Scan(&set_name, &set_data)
        if err != nil {
            panic(err)
        }

        if set_name == "user_id" {
            var ip_pre string
            var ip_render string

            ip_pre = tool.IP_preprocess(db, set_data, config.IP)[0]
            ip_render = tool.IP_parser(db, set_data, config.IP)

            data_list["user_id"] = ip_pre
            data_list["user_id_render"] = ip_render
        } else {
            data_list[set_name] = set_data
        }
    }

    return_data := make(map[string]interface{})

    if !tool.Check_acl(db, "", "", "bbs_view", config.IP) {
        data_list = map[string]string{}
        return_data["response"] = "require auth"
    }

    if other_set["legacy"] != "" {
        json_data, _ := json.Marshal(data_list)
        return string(json_data)
    } else {
        return_data["language"] = map[string]string{}
        return_data["data"] = data_list

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
}
