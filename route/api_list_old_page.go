package route

import (
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_old_page(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
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

    query := ""
    if other_set["set_type"] == "old" {
        query = tool.DB_change("select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' and not (doc_name) in (select doc_name from data_set where set_name = 'doc_type' and set_data != '') order by set_data asc limit ?, 50")
    } else {
        query = tool.DB_change("select doc_name, set_data from data_set where set_name = 'last_edit' and doc_rev = '' and not (doc_name) in (select doc_name from data_set where set_name = 'doc_type' and set_data != '') order by set_data desc limit ?, 50")
    }

    rows := tool.Query_DB(
        db,
        query,
        page_int,
    )
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var doc_name string
        var date string

        err := rows.Scan(&doc_name, &date)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{doc_name, date})
    }

    return_data := make(map[string]any)
    return_data["data"] = data_list

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
