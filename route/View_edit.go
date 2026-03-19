package route

import (
	"opennamu/route/tool"
)

func View_edit(config tool.Config, doc_name string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    raw_data := Api_w_raw(config, doc_name, "", "")
    raw_data_get := ""
    if raw_data["response"].(string) == "ok" {
        raw_data_get = raw_data["data"].(string)
    }

    form_data := `<form action="/edit/` + tool.Url_parser(doc_name) + `" method="post">
        <input class="__ON_INPUT__" type="text" name="send" placeholder="` + tool.Get_language(db, "why", true) + `">
        <hr class="main_hr">
        <textarea class="opennamu_textarea_500" id="opennamu_edit_textarea" name="content">` + tool.HTML_escape(raw_data_get) + `</textarea>
        <hr class="main_hr">
        <input class="__ON_INPUT__" type="checkbox" name="copyright_agreement">
        <hr class="main_hr">
        <button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button>
    </form>`

    out := tool.Get_template(
        db,
        config,
        doc_name,
        form_data,
        []any{ "(" + tool.Get_language(db, "edit", true) + ")" },
        [][]any{},
        map[string]string{},
    )

    return out
}
