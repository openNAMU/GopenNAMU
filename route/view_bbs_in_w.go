package route

import "opennamu/route/tool"

func View_bbs_in_w_comment(user_name string, set_id string, set_code string) string {
    return ""
}

func View_bbs_in_w(config tool.Config, set_id string, set_code string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
        return tool.Get_error_page(db, config, "auth")
    }

    bbs_name := Api_bbs_num_to_name(db, set_id)
    bbs_comment_acl := tool.Check_acl(db, set_id, "", "bbs_comment", config.IP)

    data_api := Api_bbs_w(config, set_id + "-" + set_code)
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

    if bbs_comment_acl {
        data_html += `
            <hr class="main_hr">
            <div id="opennamu_bbs_w_post_tabom"></div>
        `
    }

    data_html += `
        <hr>
        <div id="opennamu_bbs_w_post"></div>
        <script defer src="/views/main_css/js/route/topic.js` + tool.Cache_v() + `"></script>
        <script defer src="/views/main_css/js/route/bbs_w_post.js` + tool.Cache_v() + `"></script>
        <script>window.addEventListener("DOMContentLoaded", function() { opennamu_load_comment(); });</script>
    `

    if bbs_comment_acl {
        data_html += `
            <form method="post">
                <div id="opennamu_bbs_w_post_select"></div>
                ` + tool.Get_editor_ui(db, config, "", "bbs_comment", "", "") + `
            </form>
        `
    }

    Api_bbs_w_page_view_post(config, set_id, set_code)

    view_count_api := Api_bbs_w_page_view(config, set_id, set_code)
    view_count_api_data := view_count_api["data"].(int)

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