package route

import (
	"opennamu/route/tool"
)

func View_history(config tool.Config, doc_name string, set_type string, num string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok" 

    api_data := Api_list_history(config, doc_name, set_type, num)
    data_html := Get_ui_history(db, config, api_data["data"].([][]string))

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