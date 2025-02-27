package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_edit_move_all(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    if other_set["select"] == "include" {

    } else if other_set["select"] == "start" {

    } else {

    }

    return "{}"
}
