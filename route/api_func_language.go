package route

import (
	"opennamu/route/tool"
	"strings"
)

func Api_func_language(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    temp_list := strings.Split(other_set["data"], " ")

    if other_set["legacy"] != "" {
        data_list := map[string][]string{}
        data_list["data"] = []string{}

        for for_a := 0; for_a < len(temp_list); for_a++ {
            data_list["data"] = append(data_list["data"], tool.Get_language(db, temp_list[for_a], false))
        }

        json_data, _ := json.Marshal(data_list)
        return string(json_data)
    } else {
        new_data := make(map[string]any)
        new_data["response"] = "ok"

        data_list := map[string]string{}

        for for_a := 0; for_a < len(temp_list); for_a++ {
            data_list[temp_list[for_a]] = tool.Get_language(db, temp_list[for_a], false)
        }

        new_data["data"] = data_list

        json_data, _ := json.Marshal(new_data)
        return string(json_data)
    }
}
