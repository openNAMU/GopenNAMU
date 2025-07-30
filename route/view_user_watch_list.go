package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_user_watch_list(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(config.Other_set), &other_set)

	api_data := Api_user_watch_list(config, other_set["name"], other_set["num"], other_set["do_type"])
	data_html := ""

	if api_data["response"] != "ok" {
		return_data := make(map[string]any)
		return_data["response"] = "error"
		return_data["data"] = tool.Get_error_page(db, config, "auth")

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	} else {
		data_html += "<ul>"
		for _, title := range api_data["data"].([]string) {
			data_html += "<li><a href=\"/w/" + tool.Url_parser(title) + "\">" + title + "</a></li>"
		}
		data_html += "</ul>"
	}

	out := tool.Get_template(db, config, tool.Get_language(db, "watch_list", true), data_html)

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = out

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}