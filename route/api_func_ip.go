package route

import (
	"opennamu/route/tool"
)

func Api_func_ip(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    ip_data := tool.IP_parser(db, other_set["data"], config.IP)

    new_data := make(map[string]any)
    new_data["response"] = "ok"
    new_data["data"] = ip_data

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
