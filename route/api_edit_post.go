package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_edit_post_exter(config tool.Config) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_edit_post(config, other_set["doc_name"], other_set["data"])

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_edit_post(config tool.Config, doc_name string, data string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    return return_data
}