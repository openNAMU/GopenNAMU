package route

import "opennamu/route/tool"

func Api_bbs_w_page_view_exter(config tool.Config) string {
    return_data := Api_bbs_w_page_view(config)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_bbs_w_page_view(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    return return_data
}