package route

import "opennamu/route/tool"

func View_main_search_post(config tool.Config, search_type string, keyword string) string {
    if search_type == "" {
	    return tool.Get_redirect("/search/" + tool.Url_parser(keyword))
    } else if search_type == "title" {
        return tool.Get_redirect("/search_page/1/" + tool.Url_parser(keyword))
    } else {
        return tool.Get_redirect("/search_data_page/1/" + tool.Url_parser(keyword))
    }
}