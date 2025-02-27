package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_auth_post(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    ip := config.IP
    what := other_set["what"]

    tool.Do_insert_auth_history(db, ip, what)

    new_data := make(map[string]interface{})
    new_data["response"] = "ok"

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
