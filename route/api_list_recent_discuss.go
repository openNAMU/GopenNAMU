package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_recent_discuss(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    limit_int, err := strconv.Atoi(other_set["limit"])
    if err != nil {
        panic(err)
    }

    if limit_int > 50 || limit_int < 0 {
        limit_int = 50
    }

    page_int, err := strconv.Atoi(other_set["num"])
    if err != nil {
        panic(err)
    }

    if page_int > 0 {
        page_int = (page_int * limit_int) - limit_int
    } else {
        page_int = 0
    }

    set_type := other_set["set_type"]
    query := ""
    switch set_type {
    case "normal":
        query = tool.DB_change("select title, sub, date, code, stop, agree from rd order by date desc limit ?, ?")
    case "close":
        query = tool.DB_change("select title, sub, date, code, stop, agree from rd where stop = 'O' order by date desc limit ?, ?")
    default:
        query = tool.DB_change("select title, sub, date, code, stop, agree from rd where stop != 'O' order by date desc limit ?, ?")
    }

    rows := tool.Query_DB(
        db,
        query,
        page_int, limit_int,
    )
    defer rows.Close()

    data_list := [][]string{}
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var title string
        var sub string
        var date string
        var code string
        var stop string
        var agree string

        err := rows.Scan(&title, &sub, &date, &code, &stop, &agree)
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
            title,
            sub,
            date,
            code,
            stop,
            ip_pre,
            ip_render,
            id,
            agree,
        })
    }

    if other_set["legacy"] != "" {
        json_data, _ := json.Marshal(data_list)
        return string(json_data)
    } else {
        auth_name := tool.Get_user_auth(db, config.IP)
        auth_info := tool.Get_auth_group_info(db, auth_name)

        return_data := make(map[string]interface{})
        return_data["language"] = map[string]string{
            "tool":              tool.Get_language(db, "tool", false),
            "normal":            tool.Get_language(db, "normal", false),
            "close_discussion":  tool.Get_language(db, "close_discussion", false),
            "open_discussion":   tool.Get_language(db, "open_discussion", false),
            "closed":            tool.Get_language(db, "closed", false),
            "agreed_discussion": tool.Get_language(db, "agreed_discussion", false),
            "stop":              tool.Get_language(db, "stop", false),
            "admin_tool":        tool.Get_language(db, "admin_tool", false),
        }
        return_data["auth"] = auth_info
        return_data["data"] = data_list

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
}
