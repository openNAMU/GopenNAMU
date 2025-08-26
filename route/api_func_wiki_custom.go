package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_wiki_custom(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    wiki_custom_set := tool.Get_wiki_custom(db, config.IP, config.Session, config.Cookies)

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["data"] = wiki_custom_set

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}