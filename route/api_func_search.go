package route

import (
	"opennamu/route/tool"
)

func Api_func_search(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page := tool.Str_to_int(other_set["num"])
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    
    name := other_set["name"]
    query := ""
    
    if other_set["search_type"] == "title" {
        name = tool.Do_remove_spaces(name)
        query = "select title from data where replace(title, ' ', '') collate nocase like ? order by title limit ?, 50"
    } else {
        query = "select title from data where data collate nocase like ? order by title limit ?, 50"
    }

    title_list := []string{}

    rows := tool.Query_DB(
        db,
        query,
        "%" + name + "%", num,
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
