package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_acl(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    data := tool.List_acl(other_set["type"])

    return_data := make(map[string]interface{})
    return_data["response"] = "ok"
    return_data["data"] = data

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
