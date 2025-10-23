package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_user_watch_list(config tool.Config, name string, num_str string, do_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page := tool.Str_to_int(num_str)
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    ip := config.IP

    return_data := make(map[string]any)
    return_data["language"] = map[string]string{
        "watchlist": tool.Get_language(db, "watchlist", false),
        "star_doc":  tool.Get_language(db, "star_doc", false),
    }

    if ip != name && !tool.Check_acl(db, "", "", "view_user_watchlist", ip) {
        return_data["response"] = "require auth"
        return_data["data"] = []string{}
    } else {
        query := ""
        if do_type == "star_doc" {
            query = "select data from user_set where name = 'star_doc' and id = ? limit ?, 50"
        } else {
            query = "select data from user_set where name = 'watchlist' and id = ? limit ?, 50"
        }

        rows := tool.Query_DB(
            db,
            query,
            name, num,
        )
        defer rows.Close()

        data_list := []string{}

        for rows.Next() {
            var title_data string

            err := rows.Scan(&title_data)
            if err != nil {
                panic(err)
            }

            data_list = append(data_list, title_data)
        }

        return_data["response"] = "ok"
        return_data["data"] = data_list
    }

    return return_data
}

func Api_user_watch_list_exter(config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_user_watch_list(config, other_set["name"], other_set["num"], other_set["do_type"])

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
