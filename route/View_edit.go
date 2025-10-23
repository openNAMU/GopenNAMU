package route

import (
	"log"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_edit(config tool.Config, doc_name string) tool.View_result {
	db := tool.DB_connect()
	defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    raw_data := Api_w_raw(config, doc_name, "", "")
    raw_data_get := ""
    if raw_data["response"].(string) == "ok" {
        raw_data_get = raw_data["data"].(string)
    }

    log.Default().Println(raw_data)

    form_data := `<form action="/edit/` + tool.Url_parser(doc_name) + `" method="post">
        <input type="text" name="send" placeholder="` + tool.Get_language(db, "why", true) + `">
        <hr class="main_hr">
        <textarea class="opennamu_textarea_500" id="opennamu_edit_textarea" name="content">` + tool.HTML_escape(raw_data_get) + `</textarea>
        <hr class="main_hr">
        <input type="checkbox" name="copyright_agreement">
        <hr class="main_hr">
        <button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button>
    </form>`

    out := tool.Get_template(
        db,
        config,
        doc_name,
        form_data,
        "(" + tool.Get_language(db, "edit", true) + ")",
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