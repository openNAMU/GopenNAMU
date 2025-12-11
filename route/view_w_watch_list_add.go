package route

import "opennamu/route/tool"

func View_w_watch_list_add(config tool.Config, doc_name string, do_type string) tool.View_result {
	db := tool.DB_connect()
	defer tool.DB_close(db)

    name_from := false
    switch do_type {
    case "watchlist_from":
        name_from = true
        do_type = "watchlist"
    case "star_doc_from":
        name_from = true
        do_type = "star_doc"
    }

    if do_type != "watchlist" {
        do_type = "star_doc"
    }

    api_data := Api_w_watch_list_post(config, doc_name, do_type)

    return_data := make(map[string]any)

    if api_data["response"] != "ok" {
        return_data["response"] = "error"
        return_data["data"] = tool.Get_error_page(db, config, "error")
    } else {
        return_data["response"] = "ok"

        if name_from {
            return_data["data"] = tool.Get_redirect("/w/" + doc_name)
        } else if do_type == "watchlist" {
            return_data["data"] = tool.Get_redirect("/watch_list")
        } else {
            return_data["data"] = tool.Get_redirect("/star_doc")
        }
    }

    json_data, _ := json.Marshal(return_data)

    data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }
    
    return data
}