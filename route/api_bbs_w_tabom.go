package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_tabom(db *sql.DB, call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    sub_code := other_set["sub_code"]
    sub_code_parts := strings.Split(sub_code, "-")

    bbs_num := ""
    post_num := ""

    if len(sub_code_parts) > 1 {
        bbs_num = sub_code_parts[0]
        post_num = sub_code_parts[1]
    }
    
    return_data := make(map[string]interface{})

    if !tool.Check_acl(db, "", "", "bbs_comment", other_set["ip"]) {
        return_data["response"] = "require auth"
        return_data["data"] = "0"
    } else {
        stmt, err := db.Prepare(tool.DB_change("select set_data from bbs_data where set_name = 'tabom_count' and set_id = ? and set_code = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()
    
        var tabom_count string
    
        err = stmt.QueryRow(bbs_num, post_num).Scan(&tabom_count)
        if err != nil {
            if err == sql.ErrNoRows {
                tabom_count = "0"
            } else {
                panic(err)
            }
        }
    
        return_data["response"] = "ok"
        return_data["data"] = tabom_count
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}