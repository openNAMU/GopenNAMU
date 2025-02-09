package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_email(db *sql.DB, call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    err := tool.Send_email(db, other_set["ip"], other_set["who"], other_set["title"], other_set["data"])
    if err == nil {
        new_data := make(map[string]interface{})
        new_data["response"] = "ok"
    
        json_data, _ := json.Marshal(new_data)
        return string(json_data)
    }

    new_data := make(map[string]interface{})
    new_data["response"] = "err"
    new_data["data"] = err.Error()

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}