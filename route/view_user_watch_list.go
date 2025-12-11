package route

import (
	"opennamu/route/tool"
)

func View_user_watch_list(config tool.Config, num string, do_type string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    api_data := Api_user_watch_list(config, config.IP, num, do_type)
    data_html := ""

    if api_data["response"] != "ok" {
        return_data := make(map[string]any)
        return_data["response"] = "error"
        return_data["data"] = tool.Get_error_page(db, config, "auth")

        json_data, _ := json.Marshal(return_data)

        data := tool.View_result{
            HTML : return_data["data"].(string),
            JSON : string(json_data),
        }

        return data
    } else {
        data_html += "<ul>"
        for _, title := range api_data["data"].([]string) {
            data_html += "<li><a href=\"/w/" + tool.Url_parser(title) + "\">" + tool.HTML_escape(title) + "</a></li>"
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
        "",
        [][]any{
            { "user/" + tool.Url_parser(config.IP), tool.Get_language(db, "return", false) },
            { "watch_list", tool.Get_language(db, "watchlist", false) },
            { "star_doc", tool.Get_language(db, "star_doc", false) },
        },
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out
    
    json_data, _ := json.Marshal(return_data)

    data := tool.View_result{
        HTML : out,
        JSON : string(json_data),
    }

    return data
}