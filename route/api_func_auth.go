package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_auth(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    ip := config.IP
    if _, exist := other_set["ip"]; exist {
        ip = other_set["ip"]
    }

    auth_name := tool.Get_user_auth(db, ip)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["name"] = auth_name
    return_data["info"] = auth_info

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
