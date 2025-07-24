package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_list_random(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data_list := []string{}

	for i := 0; i < 50; i++ {
		title := "Test"
		tool.QueryRow_DB(
			db,
			tool.DB_change("select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1"),
			[]any{&title},
		)

		data_list = append(data_list, title)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list

    return return_data
}

func Api_list_random_exter(config tool.Config) string {
    return_data := Api_list_random(config)

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}