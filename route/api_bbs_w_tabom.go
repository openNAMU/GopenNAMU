package route

import (
	"opennamu/route/tool"
	"strings"
)

func Api_bbs_w_tabom_exter(config tool.Config) string {
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    code_split := strings.Split(other_set["sub_code"], "-")
    
    other_set["set_id"] = code_split[0]
    other_set["set_code"] = code_split[1]

    return_data := Api_bbs_w_tabom(config, other_set["set_id"], other_set["set_code"])

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_bbs_w_tabom(config tool.Config, set_id string, set_code string) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    return_data := make(map[string]string)

    if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
        return_data["response"] = "require auth"
        return_data["data"] = "0"
    } else {
        tabom_count := "0"
        tool.QueryRow_DB(
            db,
            "select set_data from bbs_data where set_name = 'tabom_count' and set_id = ? and set_code = ?",
            []any{ &tabom_count },
            set_id,
            set_code,
        )
    
        return_data["response"] = "ok"
        return_data["data"] = tabom_count
    }

    return return_data
}