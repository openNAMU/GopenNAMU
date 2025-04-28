package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_random(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    
    title := "Test"
    tool.QueryRow_DB(
        db,
        tool.DB_change("select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1"),
        []any{ &title },
    )

    new_data := map[string]string{}
    new_data["data"] = title

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
