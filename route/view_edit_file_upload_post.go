package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_edit_file_upload_post(config tool.Config) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    
    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : "",
        JSON : string(json_data),
    }

    return result_data
}