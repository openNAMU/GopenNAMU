package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_user(config tool.Config, user_name string) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    return_data := make(map[string]any)
    return_data["response"] = "ok" 

    login_menu := ""
    tool_menu := ""

    if user_name == config.IP {
        count := "0"

        tool.QueryRow_DB(
            db,
            ``,
            []any{ &count },
            config.IP,
        )

        tool_menu += `<li><a href="/alarm">` + tool.Get_language(db, "alarm", true) + "</a> (" + count + `)</li>`

        if !tool.IP_or_user(config.IP) {
            login_menu += `
                <li><a href="/logout">` + tool.Get_language(db, "logout", true) + `</a></li>
                <li><a href="/change">` + tool.Get_language(db, "user_setting", true) + `</a></li>
            `

            tool_menu += `<li><a href="/watch_list">` + tool.Get_language(db, `watchlist`, true) + `</a></li>`
            tool_menu += `<li><a href="/star_doc">` + tool.Get_language(db, `star_doc`, true) + `</a></li>`
            tool_menu += `<li><a href="/challenge">` + tool.Get_language(db, `challenge_and_level_manage`, true) + `</a></li>`
            tool_menu += `<li><a href="/acl/user:` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `user_document_acl`, true) + `</a></li>`
        } else {
            login_menu += `
                <li><a href="/login">` + tool.Get_language(db, `login`, true) + `</a></li>
                <li><a href="/register">` + tool.Get_language(db, `register`, true) + `</a></li>
                <li><a href="/change">` + tool.Get_language(db, `user_setting`, true) + `</a></li>
                <li><a href="/login/find">` + tool.Get_language(db, `password_search`, true) + `</a></li>
            `
        }

        login_menu = `<h2>` + tool.Get_language(db, `login`, true) + `</h2><ul>` + login_menu + `</ul>`
        tool_menu = `<h2>` + tool.Get_language(db, `tool`, true) + `</h2><ul>` + tool_menu + `</ul>`
    }

    admin_menu := ""

    return_data["data"] = tool.Get_template(
        db,
        config,
        tool.Get_language(db, "user_tool", true),
        `<h2>` + tool.Get_language(db, `state`, true) + `</h2>
        <div id="opennamu_get_user_info">` + tool.HTML_escape(config.IP) + `</div>
        ` + login_menu + `
        ` + tool_menu + `
        <h2>` + tool.Get_language(db, `other`, true) + `</h2>
        <ul>
            <li><a href="/record/` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `edit_record`, true) + `</a></li>
            <li><a href="/record/topic/` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `discussion_record`, true) + `</a></li>
            <li><a href="/record/bbs/` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `bbs_record`, true) + `</a></li>
            <li><a href="/record/bbs_comment/` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `bbs_comment_record`, true) + `</a></li>
            <li><a href="/topic/user:` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `user_discussion`, true) + `</a></li>
            <li><a href="/count/` + tool.Url_parser(config.IP) + `">` + tool.Get_language(db, `count`, true) + `</a></li>
        </ul>
        ` + admin_menu,
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