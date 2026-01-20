package route

import (
	"opennamu/route/tool"
)

func View_bbs_in(config tool.Config, bbs_num string, page_num string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    bbs_name := ""

    tool.QueryRow_DB(
        db,
        "select set_data from bbs_set where set_id = ? and set_name = 'bbs_name'",
        []any{ &bbs_name },
        bbs_num,
    )

    data_api := Api_bbs(config, bbs_num, page_num)
    data_api_in := data_api["data"].([]map[string]string)

    data_html := Get_bbs_list_ui(config, data_api_in, map[string]string{})

    out := tool.Get_template(
        db,
        config,
        bbs_name,
        data_html,
        []any{},
        [][]any{
            { "bbs/main", tool.Get_language(db, "return", true) },
            { "bbs/edit/" + tool.Url_parser(bbs_num), tool.Get_language(db, "add", true) },
            { "bbs/set/" + tool.Url_parser(bbs_num), tool.Get_language(db, "bbs_set", true) },
        },
    )

    return out
}