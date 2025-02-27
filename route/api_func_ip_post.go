package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_ip_post(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    ip_data := map[string]string{}

    for for_a := 1; ; for_a++ {
        for_a_str := strconv.Itoa(for_a)

        if val, ok := other_set["data_"+for_a_str]; ok {
            ip_data[val] = tool.IP_parser(db, val, config.IP)
        } else {
            break
        }
    }

    new_data := make(map[string]interface{})
    new_data["response"] = "ok"
    new_data["data"] = ip_data

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
