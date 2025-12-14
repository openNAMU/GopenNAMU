package route

import (
	"opennamu/route/tool"
)

func View_w_watch_list(config tool.Config, doc_name string, num string, do_type string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_w_watch_list(config, doc_name, num, do_type)
    data_html := ""

    if do_type != "watchlist" {
        do_type = "star_doc"
    }

    if api_data["response"] != "ok" {
        return_data := make(map[string]any)
        return_data["response"] = "error"
        return_data["data"] = tool.Get_error_page(db, config, "auth")

        json_data, _ := json.Marshal(return_data)

        data := tool.View_result{
            HTML : tool.Get_error_page(db, config, "auth"),
            JSON : string(json_data),
        }
        
        return data
    } else {
        data_html += "<ul>"
        for _, user_data := range api_data["data"].([][]string) {
            data_html += "<li>" + user_data[1] + "</li>"
        }
        
        data_html += "</ul>"
    }

    title := tool.Get_language(db, "watchlist", true)
    if do_type == "star_doc" {
        title = tool.Get_language(db, "star_doc", true)
    }

    out := tool.Get_template(
        db,
        config,
        title,
        data_html,
        "(" + doc_name + ")",
        [][]any{
            { "w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", false) },
            { "doc_watch_list/1/" + tool.Url_parser(doc_name), tool.Get_language(db, "watchlist", false) },
            { "doc_star_doc/1/" + tool.Url_parser(doc_name), tool.Get_language(db, "star_doc", false) },
        },
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)

    data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }

    return data
}