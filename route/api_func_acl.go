package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_acl(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    new_data := make(map[string]interface{})
    new_data["response"] = "ok"
    new_data["data"] = tool.Check_acl(db, other_set["name"], other_set["topic_number"], other_set["tool"], config.IP)

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
