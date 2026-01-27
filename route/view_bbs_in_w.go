package route

import (
	"database/sql"
	"opennamu/route/tool"
	"regexp"
	"strconv"
	"strings"
)

func View_bbs_in_w_comment(db *sql.DB, config tool.Config, user_name string, set_id string, set_code string) string {
    data_api := Api_bbs_w_comment(config, "around", set_id + "-" + set_code)
    data_api_in := data_api["data"].([]map[string]string)

    bbs_comment_acl := tool.Check_acl(db, set_id, "", "bbs_comment", config.IP)

    select_html := `
        <select id="opennamu_comment_select" name="comment_select">
            <option value="0">` + tool.Get_language(db, "normal", true) + `</option>
    `
    data_html := ""

    tabom_count_api := Api_bbs_w_tabom(config, set_id, set_code)
    tabom_count := tabom_count_api["data"]

    if bbs_comment_acl {
        data_html += `
            <hr class="main_hr">
            <div id="opennamu_bbs_w_post_tabom">
                <a href="javascript:void(0);" id="opennamu_tabom_button">
                    <span class="opennamu_bbs_w_post_tabom opennamu_svg opennamu_svg_tabom">&nbsp;</span>
                </a>
                <script>
                    window.addEventListener('DOMContentLoaded', function() {
                        document.getElementById('opennamu_tabom_button').addEventListener("click", function() {
                            opennamu_post_tabom("` + tool.JS_escape(set_id) + `", "` + tool.JS_escape(set_code) + `");
                        });
                    });
                </script>
                <hr class="main_hr">
                <span>` + tool.Get_language(db, "upvote", true) + `</span> <span class="opennamu_tabom_count">` + tabom_count + `</span>
            </div>
        `
    }

    data_html += "<hr>"

    var re = regexp.MustCompile(`^[0-9]+-[0-9]+-`)

    for _, v := range data_api_in {
        code_id := v["id"] + "-" + v["code"]
        code_id = re.ReplaceAllString(code_id, "")

        count := strings.Count(code_id, "-")

        select_html += `<option value="` + code_id + `">` + code_id + `</option>`

        color := "default"
        date := ""

        date += `<a href="javascript:opennamu_change_comment('` + code_id + `');">(` + tool.Get_language(db, "comment", true) + `)</a> `;
        date += `<a href="/bbs/tool/` + set_id + `/` + set_code + `/` + code_id + `">(` + tool.Get_language(db, "tool", true) + `)</a> `;
        date += v["comment_date"];

        padding_str := strconv.Itoa(20 * count)

        data_html += `<span style="padding-left: ` + padding_str + `px;"></span>`
        data_html += tool.Get_thread_ui(
            v["comment_user_id_render"],
            date,
            v["comment"],
            code_id,
            color,
            "",
            `width: calc(100% - ` + padding_str + `px);`,
            "",
        )
    }

    select_html += `</select> <a href="javascript:opennamu_return_comment();">(` + tool.Get_language(db, "return", true) + `)</a>`
    select_html += `<hr class="main_hr">`;

    if bbs_comment_acl {
        data_html += `
            <form method="post" action="/bbs/w/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">
                <div id="opennamu_bbs_w_post_select">` + select_html + `</div>
                ` + tool.Get_editor_ui(db, config, "", "bbs_comment", "", "") + `
            </form>
        `
    }

    data_html += `
        <script defer src="/views/main_css/js/route/topic.js` + tool.Cache_v() + `"></script>
        <script defer src="/views/main_css/js/route/bbs_w_post.js` + tool.Cache_v() + `"></script>
    `

    return data_html
}

func View_bbs_in_w(config tool.Config, set_id string, set_code string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
        return tool.Get_error_page(db, config, "auth")
    }

    bbs_name := Api_bbs_num_to_name(db, set_id)

    data_api := Api_bbs_w(config, set_id, set_code)
    data_api_in := data_api["data"].(map[string]string)

    if len(data_api_in) == 0 {
        return tool.Get_redirect("/bbs/main")
    }

    data_html := `
        <div class="opennamu_bbs_w_post_tab">
            <big><big><big>` + tool.HTML_escape(data_api_in["title"]) + `</big></big></big>
            <hr class="main_hr">
            ` + data_api_in["user_id_render"] + ` <span style="float: right;">` + data_api_in["date"] + `</span>
            <hr>
            <div class="opennamu_bbs_w_post_tab_content">
                ` + tool.HTML_escape(data_api_in["data"]) + `
            </div>
        </div>
    `

    Api_bbs_w_page_view_post(config, set_id, set_code)

    view_count_api := Api_bbs_w_page_view(config, set_id, set_code)
    view_count_api_data := view_count_api["data"].(int)

    data_html += View_bbs_in_w_comment(db, config, data_api_in["user_id"], set_id, set_code)

    out := tool.Get_template(
        db,
        config,
        bbs_name,
        data_html,
        []any{ "(" + tool.Get_language(db, "bbs", true) + ")", data_api_in["date"], 0, 0, view_count_api_data},
        [][]any{
            { "bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true) },
            { "bbs/edit/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "edit", true) },
            { "bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "tool", true) },
        },
    )

    return out
}