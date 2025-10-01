package route

import (
	"net/http"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_w(config tool.Config, doc_name string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

	return_data := make(map[string]any)
    return_data["response"] = "ok"

	Api_w_page_view(config)

	var render_data string
	var status int

	raw_data_api := Api_w_raw(config, doc_name, "", "")
	if raw_data_api["response"].(string) != "ok" {
		raw_data := ""

		tool.QueryRow_DB(
			db,
			`select data from other where name = "error_404"`,
			[]any{ &raw_data },
		)

		end_data := ""

		if raw_data != "" {
			end_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + raw_data + "</li></ul>"
		} else {
			end_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + tool.Get_language(db, "document_404_error", true) + "</li></ul>"
		}
		
		render_data = end_data
		status = http.StatusNotFound
	} else {
		raw_data := raw_data_api["data"].(string)
		status = http.StatusOK

		render_data_api := Api_w_render(config, doc_name, raw_data, "normal")
		render_data = render_data_api["data"]
	}

	out := tool.Get_template(
		db,
		config,
		doc_name,
		render_data,
        "",
        [][]any{
			{ "edit/" + tool.Url_parser(doc_name), tool.Get_language(db, "edit", true) },
			{ "topic/" + tool.Url_parser(doc_name), tool.Get_language(db, "discussion", true) },
			{ "history/" + tool.Url_parser(doc_name), tool.Get_language(db, "history", true) },
			{ "xref/" + tool.Url_parser(doc_name), tool.Get_language(db, "backlink", true) },
			{ "acl/" + tool.Url_parser(doc_name), tool.Get_language(db, "setting", true) },
		},
	)
	return_data["data"] = out

    json_data, _ := json.Marshal(return_data)

    data := tool.View_result{
        HTML : out,
        JSON : string(json_data),
		ST : status,
    }

    return data
}