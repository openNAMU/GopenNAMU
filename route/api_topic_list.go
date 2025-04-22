package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_topic_list(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page_int, err := strconv.Atoi(other_set["num"])
    if err != nil {
        panic(err)
    }

    if page_int > 0 {
        page_int = (page_int * 50) - 50
    } else {
        page_int = 0
    }

    rows := tool.Query_DB(
        db,
        tool.DB_change("select code, sub, stop, agree, date from rd where title = ? order by sub asc limit ?, 50"),
        other_set["name"], page_int,
    )
    defer rows.Close()

    data_list := [][]string{}
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var code string
        var sub string
        var stop string
        var agree string
        var date string

        err := rows.Scan(&code, &sub, &stop, &agree, &date)
        if err != nil {
            panic(err)
        }

        ip := ""
        id := ""
        tool.QueryRow_DB(
            db,
            tool.DB_change("select ip, id from topic where code = ? order by id + 0 desc limit 1"),
            []any{ &ip, &id },
            code,
        )

        var ip_pre string
        var ip_render string

        if _, ok := ip_parser_temp[ip]; ok {
            ip_pre = ip_parser_temp[ip][0]
            ip_render = ip_parser_temp[ip][1]
        } else {
            ip_pre = tool.IP_preprocess(db, ip, config.IP)[0]
            ip_render = tool.IP_parser(db, ip, config.IP)

            ip_parser_temp[ip] = []string{ip_pre, ip_render}
        }

        data_list = append(data_list, []string{
            code,
            sub,
            stop,
            agree,
            ip_pre,
            ip_render,
            date,
            id,
        })
    }

    return_data := make(map[string]any)
    return_data["language"] = map[string]string{
        "closed":            tool.Get_language(db, "closed", false),
        "agreed_discussion": tool.Get_language(db, "agreed_discussion", false),
        "make_new_topic":    tool.Get_language(db, "make_new_topic", false),
        "stop":              tool.Get_language(db, "stop", false),
    }
    return_data["data"] = data_list

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
