package route

import (
    "database/sql"
    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Api_func_auth_post(db *sql.DB, call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    ip := other_set["ip"]
    what := other_set["what"]

    tool.Do_insert_auth_history(db, ip, what)

    new_data := make(map[string]interface{})
    new_data["response"] = "ok"

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
