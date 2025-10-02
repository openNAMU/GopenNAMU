package route

import (
	"opennamu/route/tool"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Get_safe_send_data(data string) string {
    return tool.HTML_escape(data)
}

func View_list_recent_change(config tool.Config, set_type string, limit string, num string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := make(map[string]any)
    return_data["response"] = "ok" 

    api_data := Api_list_recent_change(config, set_type, limit, num)

    auth_name := tool.Get_user_auth(db, config.IP)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    date_heading := ""
    data_html := ""

    for_count := 0
    for _, in_data := range api_data["data"].([][]string) {
        for_count_str := strconv.Itoa(for_count)
        for_count += 1

        if in_data[6] != "" && in_data[1] == "" {
            if date_heading != "----" {
                data_html += "<h2>----</h2>"
                date_heading = "----"
            }

            data_html += tool.Get_list_ui("----", "", "", "")
            continue
        }

        doc_name := in_data[1]
        doc_name_url := tool.Url_parser(doc_name)
        rev_str := in_data[0]

        left := `<a href="/w/` + doc_name_url + `">` + tool.HTML_escape(doc_name) + `</a> `
        rev := ""

        if in_data[6] != "" {
            rev = `<span style="color: red;">r` + rev_str + `</span>`
        } else {
            rev = `r` + rev_str
        }

        rev_int := tool.Str_to_int(rev_str)
        if rev_int > 1 {
            before_rev := rev_int - 1
            before_rev_str := strconv.Itoa(before_rev)

            rev = `<a href="/diff/` + before_rev_str + `/` + rev_str + `/` + doc_name_url + `">` + rev + `</a>'`
        }

        right := ""
        right += `<span id="opennamu_list_history_` + for_count_str + `_over">`
        right += `<a id="opennamu_list_history_` + for_count_str + `" href="javascript:void(0);">`
        right += `<span class="opennamu_svg opennamu_svg_tool">&nbsp;</span></a>`
        right += `<span class="opennamu_popup_footnote" id="opennamu_list_history_` + for_count_str + `_load" style="display: none;"></span>`
        right += `</span> | `
        right += rev + " | "

        diff_size := in_data[5]
        if diff_size == "0" {
            right += `<span style="color: gray;">` + diff_size + `</span>`
        } else if strings.Contains(diff_size, "+") {
            right += `<span style="color: green;">` + diff_size + `</span>`
        } else {
            right += `<span style="color: red;">` + diff_size + `</span>`
        }

        right += " | "
        right += in_data[7] + " | "

        edit_type := "edit"
        if in_data[8] != "" {
            edit_type = in_data[8]
        }

        right += tool.Get_language(db, edit_type, true) + " | "

        time_split := strings.Split(in_data[2], " ")
        if date_heading != time_split[0] {
            data_html += "<h2>" + time_split[0] + "</h2>"
            date_heading = time_split[0]
        }

        if len(time_split) > 1 {
            right += time_split[1]
        }

        right += `<span style="display: none;" id="opennamu_history_tool_` + for_count_str + `">`

        right += `<a href="/render/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "view", true) + `</a>`
        right += ` | <a href="/raw/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "raw", true) + `</a>`
        right += ` | <a href="/revert/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "revert", true) + ` (r` + rev_str + `)</a>`

        if rev_int > 1 {
            before_rev := rev_int - 1
            before_rev_str := strconv.Itoa(before_rev)

            right += ` | <a href="/revert/` + before_rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "revert", true) + ` (r` + before_rev_str + `)</a>`
            right += ` | <a href="/diff/` + before_rev_str + `/` + rev_str +  `/` + doc_name_url + `">` + tool.Get_language(db, "compare", true) + `</a>`
        }

        right += ` | <a href="/history/` + doc_name_url + `">` + tool.Get_language(db, "history", true) + `</a>`

        if _, ok := auth_info["owner"]; ok {
            right += ` | <a href="/history_hidden/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "hide", true) + `</a>`
            right += ` | <a href="/history_delete/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "history_delete", true) + `</a>`
            right += ` | <a href="/history_send/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "send_edit", true) + `</a>`
        } else if _, ok := auth_info["hidel"]; ok {
            right += ` | <a href="/history_hidden/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "hide", true) + `</a>`
        }

        right += `</span>`

        bottom := ``
        if in_data[4] != "" {
            bottom = Get_safe_send_data(in_data[4])
        }

        data_html += tool.Get_list_ui(left, right, bottom, "")

        data_html += `<script>
            document.getElementById('opennamu_list_history_` + for_count_str + `').addEventListener("click", function() {{
                opennamu_do_footnote_popover('opennamu_list_history_` + for_count_str + `', '', 'opennamu_history_tool_` + for_count_str + `', 'open');
            }});
            document.addEventListener("click", function() {{
                opennamu_do_footnote_popover('opennamu_list_history_` + for_count_str + `', '', 'opennamu_history_tool_` + for_count_str + `', 'close');
            }});
        </script>`
    }

    return_data["data"] = tool.Get_template(
        db,
        config,
        tool.Get_language(db, "recent_change", true),
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