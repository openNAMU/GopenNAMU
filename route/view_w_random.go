package route

import (
	"opennamu/route/tool"
)

func View_w_random(config tool.Config) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_list_random(config, 1)
    api_list := api_data["data"].([]string)

    title := "FrontPage"
    if len(api_list) > 0 {
        title = api_list[0]
    }

    redirect := tool.Get_redirect("/w/" + tool.Url_parser(title))

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = redirect

    json_data, _ := json.Marshal(return_data)

    data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }

    return data
}