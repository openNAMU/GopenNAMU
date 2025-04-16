package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_recent_edit_request(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    limit_int, err := strconv.Atoi(other_set["limit"])
    if err != nil {
        panic(err)
    }

    if limit_int > 50 || limit_int < 0 {
        limit_int = 50
    }

    rows := tool.Query_DB(
        db,
        tool.DB_change("select doc_name, doc_rev, set_data from data_set where set_name = 'edit_request_doing' order by set_data desc limit ?"),
        limit_int,
    )
    defer rows.Close()

    var doc_name string
    var doc_rev string
    var date string

    data_list := [][]string{}

    for rows.Next() {
        err := rows.Scan(&doc_name, &doc_rev, &date)
        if err != nil {
            panic(err)
        }

        ip := ""
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from data_set where set_name = 'edit_request_user' and doc_rev = ? and doc_name = ?"),
            []any{ &ip },
            doc_rev, doc_name,
        )

        send := ""
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from data_set where set_name = 'edit_request_send' and doc_rev = ? and doc_name = ?"),
            []any{ &send },
            doc_rev, doc_name,
        )

        leng := ""
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from data_set where set_name = 'edit_request_leng' and doc_rev = ? and doc_name = ?"),
            []any{ &leng },
            doc_rev, doc_name,
        )

        data_list = append(data_list, []string{
            doc_name,
            doc_rev,
            date,
            tool.IP_preprocess(db, ip, config.IP)[0],
            send,
            leng,
            tool.IP_parser(db, ip, config.IP),
        })
    }

    json_data, _ := json.Marshal(data_list)
    return string(json_data)
}
