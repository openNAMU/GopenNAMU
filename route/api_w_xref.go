package route

import (
	"opennamu/route/tool"
)

func Api_w_xref(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page := tool.Str_to_int(other_set["page"])
    num := 0
    if page * 50 > 0 {
        num = page * 50 - 50
    }

    link_case_insensitive := ""
    tool.QueryRow_DB(
        db,
        "select data from other where name = 'link_case_insensitive'",
        []any{ &link_case_insensitive },
        other_set["name"],
    )

    if link_case_insensitive != "" {
        link_case_insensitive = " collate nocase"
    }

    query := ""
    if other_set["do_type"] == "1" {
        query = "select distinct link, type from back where title" + link_case_insensitive + " = ? and not type = 'no' and not type = 'nothing' order by type asc, link asc limit ?, 50"
    } else {
        query = "select distinct title, type from back where link" + link_case_insensitive + " = ? and not type = 'no' and not type = 'nothing' order by type asc, title asc limit ?, 50"
    }

    rows := tool.Query_DB(
        db,
        query,
        other_set["name"], num,
    )
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var name string
        var type_data string

        err := rows.Scan(&name, &type_data)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{name, type_data})
    }

    json_data, _ := json.Marshal(data_list)
    return string(json_data)
}
