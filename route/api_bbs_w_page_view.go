package route

import "opennamu/route/tool"

func Api_bbs_w_page_view(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    return return_data
}