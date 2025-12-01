package route

import (
	"opennamu/route/tool"
)

func Api_user_info(config tool.Config) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data := map[string]string{}

    return data
}

func Api_user_info_exter(config tool.Config) string {
    return_data := Api_user_info(config)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}