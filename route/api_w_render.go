package route

import (
	"opennamu/route/tool"
	"opennamu/route/tool/markup"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_render(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    data := markup.Get_render(db, other_set["doc_name"], other_set["data"], other_set["render_type"])

    json_data, _ := json.Marshal(data)
    return string(json_data)
}
