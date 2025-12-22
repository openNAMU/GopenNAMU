package route

import (
	"opennamu/route/tool"
)

func View_edit_file_upload(config tool.Config) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_html := `<form method="post" enctype="multipart/form-data" accept-charset="utf8">`
    data_html += `<input class="__ON_INPUT__" multiple="multiple" type="file" name="f_data[]" id="file_input">`
    data_html += `<hr class="main_hr">`
    data_html += `<input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "file_name", false) + `" name="f_name" value="">`
    data_html += `<hr class="main_hr">`
    data_html += `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", false) + `</button>`
    data_html += `</form>`

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "upload", true),
        data_html,
        "",
        [][]any{},
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : out,
        JSON : string(json_data),
    }

    return result_data
}