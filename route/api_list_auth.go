package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_auth(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    data := tool.List_auth(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["language"] = map[string]string{
        "send" : tool.Get_language(db, "send", false),
        "many_delete_help" : tool.Get_language(db, "many_delete_help", false),
    }
    return_data["data"] = data

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
