package route

import (
	"opennamu/route/tool"

	"github.com/flosch/pongo2/v6"
	jsoniter "github.com/json-iterator/go"
)

func View_render(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

	context := pongo2.Context{
		"imp" : []any{
			"test",
			tool.Get_wiki_set(db, config.IP, config.Cookies),
			tool.Get_wiki_custom(db, config.IP, config.Session, config.Cookies),
			tool.Get_wiki_css([]any{0, 0}, config.Cookies),
		},
		"data" : "test",
		"menu" : 0,
	}

	tpl, err := pongo2.FromFile(tool.Get_skin_route(db, config.IP))
	if err != nil {
		panic(err)
	}

	out, err := tpl.Execute(context)
	if err != nil {
		panic(err)
	}

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}