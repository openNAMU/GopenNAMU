package route

import (
	"opennamu/route/tool"
)

func View_list_random(config tool.Config) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_list := Api_list_random(config, 50)

    data_html := "<ul>"
    for _, title := range data_list["data"].([]string) {
        data_html += "<li><a href=\"/w/" + tool.Url_parser(title) + "\">" + tool.HTML_escape(title) + "</a></li>"
    }
    
    data_html += "</ul>"

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "random_list", true),
        data_html,
        []any{},
        [][]any{
            { "other", tool.Get_language(db, "return", true) },
        },
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : out,
        JSON : string(json_data),
    }

    return result_data
}