package route

import (
    "database/sql"
    "encoding/json"
    "opennamu/route/tool"
)

func Api_list_auth(db *sql.DB, call_arg []string) string {
    data := tool.List_auth(db)

    return_data := make(map[string]interface{})
    return_data["response"] = "ok"
    return_data["language"] = map[string]string{
        "send":             tool.Get_language(db, "send", false),
        "many_delete_help": tool.Get_language(db, "many_delete_help", false),
    }
    return_data["data"] = data

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
