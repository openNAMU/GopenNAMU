package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

type Api_list_recent_change_T struct {
    full map[string]any
    legacy [][]string
}

func Api_list_recent_change_exter(config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_list_recent_change_call(config, other_set["type"], other_set["limit"], other_set["num"])

    var json_data []byte

    if other_set["legacy"] != "" {
        json_data, _ = json.Marshal(return_data.legacy)
    } else {
        json_data, _ = json.Marshal(return_data.full)
    }
    
    return string(json_data)
}

func Api_list_recent_change(config tool.Config, set_type string, limit string, num string) map[string]any {
    return Api_list_recent_change_call(config, set_type, limit, num).full
}

func Api_list_recent_change_call(config tool.Config, set_type string, limit string, num string) Api_list_recent_change_T {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if set_type == "edit" {
        set_type = ""
    }

    limit_int := tool.Str_to_int(limit)
    if limit_int > 50 || limit_int < 0 {
        limit_int = 50
    }

    page_int := tool.Str_to_int(num)
    if page_int > 0 {
        page_int = (page_int * limit_int) - limit_int
    } else {
        page_int = 0
    }

    rows := tool.Query_DB(
        db,
        tool.DB_change("select id, title from rc where type = ? order by date desc limit ?, ?"),
        set_type, page_int, limit_int,
    )
    defer rows.Close()

    data_list := [][]string{}

    admin_auth := tool.Check_acl(db, "", "", "hidel_auth", config.IP)
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var id string
        var title string

        err := rows.Scan(&id, &title)
        if err != nil {
            panic(err)
        }

        date := ""
        ip := ""
        send := ""
        leng := ""
        hide := ""
        type_data := ""
        tool.QueryRow_DB(
            db,
            tool.DB_change("select date, ip, send, leng, hide, type from history where id = ? and title = ?"),
            []any{ &date, &ip, &send, &leng, &hide, &type_data },
            id, title,
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

        if hide == "" || admin_auth {
            data_list = append(data_list, []string{
                id,
                title,
                date,
                ip_pre,
                send,
                leng,
                hide,
                ip_render,
                type_data,
            })
        } else {
            data_list = append(data_list, []string{"", "", "", "", "", "", hide, "", ""})
        }
    }

    EOL_data := Api_list_recent_change_T{}

    EOL_data.legacy = data_list

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    EOL_data.full = return_data

    return EOL_data
}
