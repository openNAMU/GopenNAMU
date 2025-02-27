package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_user_setting_editor_delete(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    ip := config.IP
    if !tool.IP_or_user(ip) {
        tool.Exec_DB(
            db,
            "delete from user_set where id = ? and name = 'user_editor_top' and data = ?",
            ip, other_set["data"],
        )

        return_data := make(map[string]interface{})
        return_data["response"] = "ok"
        return_data["language"] = map[string]string{
            "delete": tool.Get_language(db, "delete", false),
        }

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    } else {
        return_data := make(map[string]interface{})
        return_data["response"] = "require auth"
        return_data["language"] = map[string]string{
            "authority_error": tool.Get_language(db, "authority_error", false),
        }

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
}
