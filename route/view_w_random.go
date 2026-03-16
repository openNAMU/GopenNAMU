package route

import (
	"opennamu/route/tool"
)

func View_w_random(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_list_random(config, 1)
    api_list := api_data["data"].([]string)

    title := "FrontPage"
    if len(api_list) > 0 {
        title = api_list[0]
    }

    redirect := tool.Get_redirect("/w/" + tool.Url_parser(title))

    return redirect
}