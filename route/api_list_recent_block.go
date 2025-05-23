package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_recent_block(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
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

    // private 공개 안되도록 조심할 것
    var rows *sql.Rows
    switch other_set["set_type"] {
    case "all":
        query := ""
        if other_set["why"] != "" {
            query = tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band != 'private' and why like ? order by today desc limit ?, 50")
        } else {
            query = tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band != 'private' order by today desc limit ?, 50")
        }

        if other_set["why"] != "" {
            rows = tool.Query_DB(
                db,
                query,
                other_set["why"] + "%", page_int,
            )
        } else {
            rows = tool.Query_DB(
                db,
                query,
                page_int,
            )
        }
    case "ongoing":
        rows = tool.Query_DB(
            db,
            tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where ongoing = '1' and band != 'private' order by end desc limit ?, 50"),
            page_int,
        )
    case "regex":
        rows = tool.Query_DB(
            db,
            tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'regex' order by today desc limit ?, 50"),
            page_int,
        )
    case "private":
        rows = tool.Query_DB(
            db,
            tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'private' order by today desc limit ?, 50"),
            page_int,
        )
    case "user":
        rows = tool.Query_DB(
            db,
            tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where block = ? and band != 'private' order by today desc limit ?, 50"),
            other_set["user_name"], page_int,
        )
    case "cidr":
        rows = tool.Query_DB(
            db,
            tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'cidr' order by today desc limit ?, 50"),
            page_int,
        )
    default:
        rows = tool.Query_DB(
            db,
            tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where blocker = ? and band != 'private' order by today desc limit ?, 50"),
            other_set["user_name"], page_int,
        )
    }
    defer rows.Close()

    data_list := [][]string{}
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var why string
        var block string
        var blocker string
        var end string
        var today string
        var band string
        var ongoing string

        err := rows.Scan(&why, &block, &blocker, &end, &today, &band, &ongoing)
        if err != nil {
            panic(err)
        }

        var ip_pre_blocker string
        var ip_render_blocker string

        if _, ok := ip_parser_temp[blocker]; ok {
            ip_pre_blocker = ip_parser_temp[blocker][0]
            ip_render_blocker = ip_parser_temp[blocker][1]
        } else {
            ip_pre_blocker = tool.IP_preprocess(db, blocker, config.IP)[0]
            ip_render_blocker = tool.IP_parser(db, blocker, config.IP)

            ip_parser_temp[blocker] = []string{ip_pre_blocker, ip_render_blocker}
        }

        var ip_pre_block string
        var ip_render_block string

        if band == "" {
            if _, ok := ip_parser_temp[block]; ok {
                ip_pre_block = ip_parser_temp[block][0]
                ip_render_block = ip_parser_temp[block][1]
            } else {
                ip_pre_block = tool.IP_preprocess(db, block, config.IP)[0]
                ip_render_block = tool.IP_parser(db, block, config.IP)

                ip_parser_temp[block] = []string{ip_pre_block, ip_render_block}
            }
        } else {
            ip_pre_block = block
            ip_render_block = block
        }

        data_list = append(data_list, []string{
            why,
            ip_pre_block,
            ip_render_block,
            ip_pre_blocker,
            ip_render_blocker,
            end,
            today,
            band,
            ongoing,
        })
    }

    if other_set["set_type"] == "private" {
        if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
            data_list = [][]string{}
        }
    }

    return_data := make(map[string]any)
    return_data["language"] = map[string]string{
        "all":         tool.Get_language(db, "all", false),
        "regex":       tool.Get_language(db, "regex", false),
        "cidr":        tool.Get_language(db, "cidr", false),
        "private":     tool.Get_language(db, "private", false),
        "in_progress": tool.Get_language(db, "in_progress", false),
        "admin":       tool.Get_language(db, "admin", false),
        "blocked":     tool.Get_language(db, "blocked", false),
        "limitless":   tool.Get_language(db, "limitless", false),
        "release":     tool.Get_language(db, "release", false),
        "start":       tool.Get_language(db, "start", false),
        "end":         tool.Get_language(db, "end", false),
        "ban":         tool.Get_language(db, "ban", false),
        "why":         tool.Get_language(db, "why", false),
    }
    return_data["data"] = data_list

    auth_name := tool.Get_user_auth(db, config.IP)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    return_data["auth"] = auth_info

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
