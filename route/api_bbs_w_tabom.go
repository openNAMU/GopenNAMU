package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_tabom(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    sub_code := other_set["sub_code"]
    sub_code_parts := strings.Split(sub_code, "-")

    bbs_num := ""
    post_num := ""

    if len(sub_code_parts) > 1 {
        bbs_num = sub_code_parts[0]
        post_num = sub_code_parts[1]
    }
    
    return_data := make(map[string]any)

    if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
        return_data["response"] = "require auth"
        return_data["data"] = "0"
    } else {
        tabom_count := "0"
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from bbs_data where set_name = 'tabom_count' and set_id = ? and set_code = ?"),
            []any{ &tabom_count },
            bbs_num, post_num,
        )
    
        return_data["response"] = "ok"
        return_data["data"] = tabom_count
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}