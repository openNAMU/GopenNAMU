package route

import (
    "database/sql"

    jsoniter "github.com/json-iterator/go"
)

func Api_edit_move_all(db *sql.DB, call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    if other_set["select"] == "include" {

    } else if other_set["select"] == "start" {

    } else {

    }

    return "{}"
}
