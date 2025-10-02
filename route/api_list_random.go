package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_list_random(config tool.Config, list_count int) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_list := []string{}

    rows := tool.Query_DB(
        db,
        tool.DB_change("select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit ?"),
        list_count,
    )

    for rows.Next() {
        var title string

        err := rows.Scan(&title)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, title)
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list

    return return_data
}

func Api_list_random_exter(config tool.Config) string {
    return_data := Api_list_random(config, 50)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}