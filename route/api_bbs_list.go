package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_bbs_list_exter(config tool.Config) string {
    return_data := Api_bbs_list(config)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func bbs_list(db *sql.DB) map[string]string {
    rows := tool.Query_DB(
        db,
        "select set_data, set_id from bbs_set where set_name = 'bbs_name'",
    )
    defer rows.Close()

    data_list := map[string]string{}

    for rows.Next() {
        var name string
        var id string

        err := rows.Scan(&name, &id)
        if err != nil {
            panic(err)
        }

        data_list[name] = id
    }

    return data_list
}

func Api_bbs_list(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_list := bbs_list(db)
    data_list_sub := map[string][]string{}

    for k, v := range data_list {
        bbs_type := ""
        tool.QueryRow_DB(
            db,
            "select set_data from bbs_set where set_name = 'bbs_type' and set_id = ?",
            []any{ &bbs_type },
            v,
        )

        bbs_date := ""
        tool.QueryRow_DB(
            db,
            "select set_data from bbs_data where set_id = ? and set_name = 'date' order by set_code + 0 desc limit 1",
            []any{ &bbs_date },
            v,
        )

        data_list_sub[k] = []string{v, bbs_type, bbs_date}
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_list_sub

    return return_data
}
