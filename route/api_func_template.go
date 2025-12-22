package route

import "opennamu/route/tool"

func Api_func_template(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]any{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return tool.Get_template(
        db,
        config,
        other_set["name"].(string),
        other_set["data"].(string),
        other_set["sub"].(string),
        other_set["menu"].([][]any),
    )
}