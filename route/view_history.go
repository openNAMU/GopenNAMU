package route

import (
	"opennamu/route/tool"
)

func View_history(config tool.Config, doc_name string, set_type string, num string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok" 

    if set_type == "" {
        set_type = "normal"
    }

    data_html := ""

    menu_option := []string{ "normal", "edit", "move", "delete", "revert", "r1", "setting" }
    for _, option := range menu_option {
        label := tool.Get_language(db, option, true)
        data_html += `<a href="/recent_change/1/` + option + `">(` + label + `)</a> `
    }

    api_data := Api_list_history(config, doc_name, set_type, num)
    api_data_list := api_data["data"].([][]string)

    data_html += Get_ui_history(db, config, api_data_list)
    data_html += tool.Get_page_control(
        db,
        tool.Str_to_int(num),
        len(api_data_list),
        50,
        "/history_page/{}/" + set_type + "/" + tool.Url_parser(doc_name),
    )

    return_data["data"] = tool.Get_template(
        db,
        config,
        doc_name,
        data_html,
        "(" + tool.Get_language(db, "history", true) + ")",
        [][]any{},
    )

    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }

    return result_data
}