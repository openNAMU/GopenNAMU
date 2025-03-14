package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_old_page(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page_int, err := strconv.Atoi(other_set["num"])
    if err != nil {
        panic(err)
    }

    if page_int > 0 {
        page_int = (page_int * 50) - 50
    } else {
        page_int = 0
    }

    var stmt *sql.Stmt

    if other_set["set_type"] == "old" {
        stmt, err = db.Prepare(tool.DB_change("select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' and not (doc_name) in (select doc_name from data_set where set_name = 'doc_type' and set_data != '') order by set_data asc limit ?, 50"))
    } else {
        stmt, err = db.Prepare(tool.DB_change("select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' and not (doc_name) in (select doc_name from data_set where set_name = 'doc_type' and set_data != '') order by set_data desc limit ?, 50"))
    }

    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    rows, err := stmt.Query(page_int)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var doc_name string
        var date string

        err := rows.Scan(&doc_name, &date)
        if err != nil {
            panic(err)
        }

        stmt, err = db.Prepare(tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'doc_type'"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        data_list = append(data_list, []string{doc_name, date})
    }

    return_data := make(map[string]interface{})
    return_data["data"] = data_list

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
