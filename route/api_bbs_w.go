package route

import (
	"opennamu/route/tool"
	"strings"
)

func Api_bbs_w_exter(config tool.Config) string {
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_bbs_w(config, other_set["sub_code"])

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_bbs_w(config tool.Config, sub_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    sub_code_parts := strings.Split(sub_code, "-")

    bbs_num := ""
    post_num := ""

    if len(sub_code_parts) > 1 {
        bbs_num = sub_code_parts[0]
        post_num = sub_code_parts[1]
    }

    rows := tool.Query_DB(
        db,
        "select set_name, set_data from bbs_data where set_id = ? and set_code = ?",
        bbs_num, post_num,
    )
    defer rows.Close()
    
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

    return_data := make(map[string]any)

    if !tool.Check_acl(db, "", "", "bbs_view", config.IP) {
        return_data["response"] = "require auth"
        return_data["data"] = map[string]string{}
    } else {
        return_data["response"] = "ok"
        return_data["data"] = data_list
    }

    return return_data
}
