package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_watch_list(config tool.Config, name string, num_str string, do_type string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    page := tool.Str_to_int(num_str)
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    return_data := make(map[string]any)
    return_data["language"] = map[string]string{
        "watchlist": tool.Get_language(db, "watchlist", false),
        "star_doc":  tool.Get_language(db, "star_doc", false),
    }

    if !tool.Check_acl(db, "", "", "doc_watch_list_view", config.IP) {
        return_data["response"] = "require auth"
        return_data["data"] = []string{}
    } else {
        query := ""
        if do_type == "star_doc" {
            query = tool.DB_change("select id from user_set where name = 'star_doc' and data = ? limit ?, 50")
        } else {
            query = tool.DB_change("select id from user_set where name = 'watchlist' and data = ? limit ?, 50")
        }

        rows := tool.Query_DB(
            db,
            query,
            name,
            num,
        )
        defer rows.Close()

        data_list := [][]string{}
        ip_parser_temp := map[string][]string{}

        for rows.Next() {
            var user_name string

            err := rows.Scan(&user_name)
            if err != nil {
                panic(err)
            }

            var ip_pre string
            var ip_render string

            if _, ok := ip_parser_temp[user_name]; ok {
                ip_pre = ip_parser_temp[user_name][0]
                ip_render = ip_parser_temp[user_name][1]
            } else {
                ip_pre = tool.IP_preprocess(db, user_name, config.IP)[0]
                ip_render = tool.IP_parser(db, user_name, config.IP)

                ip_parser_temp[user_name] = []string{ip_pre, ip_render}
            }

            data_list = append(data_list, []string{ip_pre, ip_render})
        }

        return_data["response"] = "ok"
        return_data["data"] = data_list
    }

    return return_data
}

func Api_w_watch_list_exter(config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_w_watch_list(config, other_set["name"], other_set["num"], other_set["do_type"])

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
