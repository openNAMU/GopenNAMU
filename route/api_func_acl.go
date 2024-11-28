package route

import (
    "database/sql"
    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Api_func_acl(db *sql.DB, call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    new_data := make(map[string]interface{})
    new_data["response"] = "ok"
    new_data["data"] = tool.Check_acl(db, other_set["name"], other_set["topic_number"], other_set["tool"], other_set["ip"])

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
