package route

import (
	"log"
	"opennamu/route/tool"
	"strconv"
)

func View_bbs_main(config tool.Config, page string) tool.View_result {
	db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok" 

    bbs_list_api_data := Api_bbs_list(config)
    log.Default().Println(bbs_list_api_data)

    bbs_id_to_name := map[string]string{}

    data_html := "<ul>"
    for bbs_name, in_data := range bbs_list_api_data["data"].(map[string][]string) {
        bbs_id := in_data[0]
        bbs_type := in_data[1]
        bbs_date := in_data[2]

        bbs_id_to_name[bbs_id] = bbs_name

        data_html += "<li>"
        data_html += "<a href=\"/bbs/in/" + tool.Url_parser(bbs_id) + "\">"
        data_html += tool.HTML_escape(bbs_name)

        if bbs_type == "comment" {
            data_html += " (" + tool.Get_language(db, "comment_base", false) + ")"
        } else {
            data_html += " (" + tool.Get_language(db, "thread_base", false) + ")"
        }

        if bbs_date != "" {
            data_html += " (" + bbs_date + ")"
        }

        data_html += "</a></li>"
    }

    data_html += "</ul>"

    bbs_api_data := Api_bbs(config, "", page)
    
    count := 0
    for _, in_data := range bbs_api_data["data"].([]map[string]string) {
        count_str := strconv.Itoa(count)
        count += 1

        bbs_title := in_data["title"]
        bbs_id := in_data["set_id"]
        bbs_code := in_data["set_code"]
        bbs_name := bbs_id_to_name[bbs_id]
        bbs_date := in_data["date"]
        bbs_user_id := in_data["user_id_render"]

        right := ""
        right += `<a href="/bbs/w/` + bbs_id + `/` + bbs_code + `">` + tool.HTML_escape(bbs_title) + `</a>`

        left := ""
        left += `<span id="opennamu_bbs_comment_` + count_str + `"></span>`
        left += `<a href="/bbs/in/` + bbs_id + `">` + bbs_name + `</a> | `
        left += bbs_user_id + " | "
        left += bbs_date

        data_html += tool.Get_list_ui(left, right, "", "")
    }

    return_data["data"] = tool.Get_template(
        db,
        config,
        tool.Get_language(db, "bbs_main", true),
        data_html,
        "",
        [][]any{},
    )

    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }

    return result_data
}