package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_history(config tool.Config) string {
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

    var rows *sql.Rows

    if other_set["set_type"] == "edit" {
        other_set["set_type"] = ""
    }

    if other_set["set_type"] == "normal" {
        rows = tool.Query_DB(
            db,
            tool.DB_change("select id, title, date, ip, send, leng, hide, type from history where title = ? order by id + 0 desc limit ?, 50"),
            other_set["doc_name"], page_int,
        )
    } else {
        rows = tool.Query_DB(
            db,
            tool.DB_change("select id, title, date, ip, send, leng, hide, type from history where title = ? and type = ? order by id + 0 desc limit ?, 50"),
            other_set["doc_name"], other_set["set_type"], page_int,
        )
    }
    defer rows.Close()

    data_list := [][]string{}

    admin_auth := tool.Check_acl(db, "", "", "hidel_auth", config.IP)
    ip_parser_temp := map[string][]string{}

    for rows.Next() {
        var id string
        var title string
        var date string
        var ip string
        var send string
        var leng string
        var hide string
        var type_data string

        err := rows.Scan(&id, &title, &date, &ip, &send, &leng, &hide, &type_data)
        if err != nil {
            panic(err)
        }

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

    auth_name := tool.Get_user_auth(db, config.IP)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    return_data := make(map[string]any)
    return_data["language"] = map[string]string{
        "tool":           tool.Get_language(db, "tool", false),
        "normal":         tool.Get_language(db, "normal", false),
        "edit":           tool.Get_language(db, "edit", false),
        "move":           tool.Get_language(db, "move", false),
        "delete":         tool.Get_language(db, "delete", false),
        "revert":         tool.Get_language(db, "revert", false),
        "new_doc":        tool.Get_language(db, "new_doc", false),
        "edit_request":   tool.Get_language(db, "edit_request", false),
        "user_document":  tool.Get_language(db, "user_document", false),
        "raw":            tool.Get_language(db, "raw", false),
        "compare":        tool.Get_language(db, "compare", false),
        "history":        tool.Get_language(db, "history", false),
        "hide":           tool.Get_language(db, "hide", false),
        "history_delete": tool.Get_language(db, "history_delete", false),
        "send_edit":      tool.Get_language(db, "send_edit", false),
        "file":           tool.Get_language(db, "file", false),
        "category":       tool.Get_language(db, "category", false),
        "setting":        tool.Get_language(db, "setting", false),
        "remove_hidden":  tool.Get_language(db, "remove_hidden", false),
    }
    return_data["data"] = data_list
    return_data["auth"] = auth_info

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
