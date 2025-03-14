package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_page_view(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    pv_continue := tool.Get_setting(db, "not_use_view_count", "")
    if len(pv_continue) == 0 || pv_continue[0][0] == "" {
        // 전체 조회수
        stmt, err := db.Prepare(tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ''"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        var view_count string

        err = stmt.QueryRow(other_set["doc_name"]).Scan(&view_count)
        if err != nil {
            if err == sql.ErrNoRows {
                view_count = "0"
            } else {
                panic(err)
            }
        }

        if view_count == "0" {
            tool.Exec_DB(
                db,
                "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'view_count', '1')",
                other_set["doc_name"],
            )
        } else {
            view_count_int, _ := strconv.Atoi(view_count)

            tool.Exec_DB(
                db,
                "update data_set set set_data = ? where doc_name = ? and set_name = 'view_count' and doc_rev = ''",
                view_count_int + 1, other_set["doc_name"],
            )
        }

        // 월간 조회수
        stmt, err = db.Prepare(tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        now_date := tool.Get_month()

        err = stmt.QueryRow(other_set["doc_name"], now_date).Scan(&view_count)
        if err != nil {
            if err == sql.ErrNoRows {
                view_count = "0"
            } else {
                panic(err)
            }
        }

        if view_count == "0" {
            tool.Exec_DB(
                db,
                "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'view_count', '1')",
                other_set["doc_name"], now_date,
            )
        } else {
            view_count_int, _ := strconv.Atoi(view_count)

            tool.Exec_DB(
                db,
                "update data_set set set_data = ? where doc_name = ? and set_name = 'view_count' and doc_rev = ?",
                view_count_int + 1, other_set["doc_name"], now_date,
            )
        }

        // 하루 조회수
        stmt, err = db.Prepare(tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        now_date = tool.Get_date()

        err = stmt.QueryRow(other_set["doc_name"], now_date).Scan(&view_count)
        if err != nil {
            if err == sql.ErrNoRows {
                view_count = "0"
            } else {
                panic(err)
            }
        }

        if view_count == "0" {
            tool.Exec_DB(
                db,
                "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'view_count', '1')",
                other_set["doc_name"], now_date,
            )
        } else {
            view_count_int, _ := strconv.Atoi(view_count)

            tool.Exec_DB(
                db,
                "update data_set set set_data = ? where doc_name = ? and set_name = 'view_count' and doc_rev = ?",
                view_count_int + 1, other_set["doc_name"], now_date,
            )
        }
    }

    return_data := make(map[string]interface{})
    return_data["response"] = "ok"

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
