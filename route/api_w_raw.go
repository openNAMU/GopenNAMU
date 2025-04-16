package route

import (
	"database/sql"

	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_raw(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    new_data := make(map[string]any)

    if !tool.Check_acl(db, other_set["name"], "", "render", config.IP) {
        new_data["response"] = "require auth"
    } else if other_set["exist_check"] != "" {
        title := ""
        exist := tool.QueryRow_DB(
            db,
            tool.DB_change("select title from data where title = ?"),
            []any{ &title },
            other_set["name"],
        )

        if !exist {
            new_data["exist"] = false
        } else {
            new_data["exist"] = true
        }

        new_data["response"] = "ok"
    } else {        
        exist := false

        data := ""
        hide := ""
        if other_set["rev"] != "" {
            exist = tool.QueryRow_DB(
                db,
                tool.DB_change("select data, hide from history where title = ? and id = ?"),
                []any{ &data, &hide },
                other_set["name"], other_set["rev"],
            )
        } else {
            exist = tool.QueryRow_DB(
                db,
                tool.DB_change("select data from data where title = ?"),
                []any{ &data },
                other_set["name"],
            )
        }

        if !exist {
            new_data["response"] = "not exist"
        } else {
            check_pass := false
            if hide != "" {
                if tool.Check_acl(db, "", "", "hidel_auth", config.IP) {
                    check_pass = true
                } else {
                    new_data["response"] = "require auth"
                }
            } else {
                check_pass = true
            }

            if check_pass {
                new_data["title"] = other_set["name"]
                new_data["data"] = data

                new_data["response"] = "ok"
            }
        }
    }

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
