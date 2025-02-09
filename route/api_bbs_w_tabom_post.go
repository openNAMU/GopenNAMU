package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_tabom_post(db *sql.DB, call_arg []string) string {
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
    } else {
        stmt1, err := db.Prepare(tool.DB_change("select set_data from bbs_data where set_name = 'tabom_list' and set_data = ? and set_id = ? and set_code = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt1.Close()

        not_exist := false
        var no_data string

        err = stmt1.QueryRow(other_set["ip"], bbs_num, post_num).Scan(&no_data)
        if err != nil {
            if err == sql.ErrNoRows {
                not_exist = true
            } else {
                panic(err)
            }
        }

        if not_exist {
            return_data["response"] = "ok"

            stmt2, err := db.Prepare(tool.DB_change("select set_data from bbs_data where set_name = 'tabom_count' and set_id = ? and set_code = ?"))
            if err != nil {
                panic(err)
            }
            defer stmt2.Close()
        
            var tabom_count string
        
            err = stmt2.QueryRow(bbs_num, post_num).Scan(&tabom_count)
            if err != nil {
                if err == sql.ErrNoRows {
                    tool.Exec_DB(
                        db,
                        "insert into bbs_data (set_name, set_data, set_id, set_code) values ('tabom_count', ?, ?, ?)",
                        tabom_count, bbs_num, post_num,
                    )
                } else {
                    panic(err)
                }
            }

            tabom_count_int, _ := strconv.Atoi(tabom_count)
            tabom_count_int += 1

            tabom_count_str := strconv.Itoa(tabom_count_int)

            tool.Exec_DB(
                db,
                "update bbs_data set set_data = ? where set_name = 'tabom_count' and set_id = ? and set_code = ?",
                tabom_count_str, bbs_num, post_num,
            )
            tool.Exec_DB(
                db,
                "insert into bbs_data (set_name, set_data, set_id, set_code) values ('tabom_list', ?, ?, ?)",
                other_set["ip"], bbs_num, post_num,
            )
        } else {
            return_data["response"] = "same user exist"
        }
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}