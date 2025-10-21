package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_edit(config tool.Config, doc_name string) tool.View_result {
	db := tool.DB_connect()
	defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "edit", true),
        "",
        "",
        [][]any{},
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)
    
    result_data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }

    return result_data
}