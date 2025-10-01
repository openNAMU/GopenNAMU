package route

import (
	"opennamu/route/tool"
	"opennamu/route/tool/markup"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_render_exter(config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    doc_name := other_set["doc_name"]
    raw_data := other_set["raw_data"]
    render_type := other_set["render_type"]

    return_data := Api_w_render(config, doc_name, raw_data, render_type)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_w_render(config tool.Config, doc_name string, raw_data string, render_type string) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data := markup.Get_render(db, doc_name, raw_data, render_type)

    return data
}
