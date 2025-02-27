package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_sha224(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    hash_str := tool.Sha224(other_set["data"])

    return_data := make(map[string]interface{})
    return_data["response"] = "ok"
    return_data["data"] = hash_str

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
