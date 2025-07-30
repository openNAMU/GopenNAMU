package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_list_random(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

	data_list := Api_list_random(config)

	data_html := "<ul>"
	for _, title := range data_list["data"].([]string) {
		data_html += "<li><a href=\"/w/" + tool.Url_parser(title) + "\">" + title + "</a></li>"
	}
	data_html += "</ul>"

	out := tool.Get_template(db, config, tool.Get_language(db, "random_list", true), data_html)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}