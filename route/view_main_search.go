package route

import "opennamu/route/tool"

func View_main_search(config tool.Config, keyword string, num string, search_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

	data_api := Api_func_search(config, keyword, num, search_type)
	data_api_in := data_api["data"].([]string)

    if keyword == "" {
        return tool.Get_redirect("/")
    }
}