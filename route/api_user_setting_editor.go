package route

import (
       "database/sql"
    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Api_user_setting_editor(db *sql.DB, call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    ip := other_set["ip"]
    if !tool.IP_or_user(ip) {
        stmt, err := db.Prepare(tool.DB_change("select data from user_set where id = ? and name = 'user_editor_top'"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        rows, err := stmt.Query(ip)
        if err != nil {
            panic(err)
        }
        defer rows.Close()

        data_list := []string{}

        for rows.Next() {
            var data string

            err := rows.Scan(&data)
            if err != nil {
                panic(err)
            }

            data_list = append(data_list, data)
        }

        return_data := make(map[string]interface{})
        return_data["response"] = "ok"
        return_data["data"] = data_list

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
