package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_search(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page, _ := strconv.Atoi(other_set["num"])
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    query := ""
    if other_set["search_type"] == "title" {
        query = tool.DB_change("select title from data where title collate nocase like ? order by title limit ?, 50")
    } else {
        query = tool.DB_change("select title from data where data collate nocase like ? order by title limit ?, 50")
    }

    title_list := []string{}

    rows := tool.Query_DB(
        db,
        query,
        "%" + other_set["name"] + "%", num,
    )
    defer rows.Close()

    for rows.Next() {
        var title string

        err := rows.Scan(&title)
        if err != nil {
            panic(err)
        }

        title_list = append(title_list, title)
    }

    json_data, _ := json.Marshal(title_list)
    return string(json_data)
}
