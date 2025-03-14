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
        stmt, err := db.Prepare(tool.DB_change("select title from data where title = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        var title string

        err = stmt.QueryRow(other_set["name"]).Scan(&title)
        if err != nil {
            if err == sql.ErrNoRows {
                new_data["exist"] = false
            } else {
                panic(err)
            }
        } else {
            new_data["exist"] = true
        }

        new_data["response"] = "ok"
    } else {
        var data string
        hide := ""

        var stmt *sql.Stmt
        var err error

        if other_set["rev"] != "" {
            stmt, err = db.Prepare(tool.DB_change("select data, hide from history where title = ? and id = ?"))
            if err != nil {
                panic(err)
            }
            defer stmt.Close()

            err = stmt.QueryRow(other_set["name"], other_set["rev"]).Scan(&data, &hide)
        } else {
            stmt, err = db.Prepare(tool.DB_change("select data from data where title = ?"))
            if err != nil {
                panic(err)
            }
            defer stmt.Close()

            err = stmt.QueryRow(other_set["name"]).Scan(&data)
        }

        if err != nil {
            if err == sql.ErrNoRows {
                new_data["response"] = "not exist"
            } else {
                panic(err)
            }
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
