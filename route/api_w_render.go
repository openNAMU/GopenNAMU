package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_render(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    data := tool.Get_render(db, other_set["doc_name"], other_set["data"], other_set["render_type"])

    json_data, _ := json.Marshal(data)
    return string(json_data)
}
