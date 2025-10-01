package route

import (
	"encoding/json"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_page_view_exter(config tool.Config) string {
    return_data := Api_w_page_view(config)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_w_page_view(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    pv_continue := tool.Get_setting(db, "not_use_view_count", "")
    if len(pv_continue) == 0 || pv_continue[0][0] == "" {
        // 전체 조회수
        view_count := "0"
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ''"),
            []any{ &view_count },
            other_set["doc_name"],
        )

        if view_count == "0" {
            tool.Exec_DB(
                db,
                "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'view_count', '1')",
                other_set["doc_name"],
            )
        } else {
            view_count_int := tool.Str_to_int(view_count)

            tool.Exec_DB(
                db,
                "update data_set set set_data = ? where doc_name = ? and set_name = 'view_count' and doc_rev = ''",
                view_count_int + 1, other_set["doc_name"],
            )
        }

        // 월간 조회수
        now_date := tool.Get_month()
        view_count = "0"
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ?"),
            []any{ &view_count },
            other_set["doc_name"], now_date,
        )

        if view_count == "0" {
            tool.Exec_DB(
                db,
                "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'view_count', '1')",
                other_set["doc_name"], now_date,
            )
        } else {
            view_count_int := tool.Str_to_int(view_count)

            tool.Exec_DB(
                db,
                "update data_set set set_data = ? where doc_name = ? and set_name = 'view_count' and doc_rev = ?",
                view_count_int + 1, other_set["doc_name"], now_date,
            )
        }

        // 하루 조회수
        now_date = tool.Get_date()
        view_count = "0"
        tool.QueryRow_DB(
            db,
            tool.DB_change("select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ?"),
            []any{ &view_count },
            other_set["doc_name"], now_date,
        )

        if view_count == "0" {
            tool.Exec_DB(
                db,
                "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'view_count', '1')",
                other_set["doc_name"], now_date,
            )
        } else {
            view_count_int := tool.Str_to_int(view_count)

            tool.Exec_DB(
                db,
                "update data_set set set_data = ? where doc_name = ? and set_name = 'view_count' and doc_rev = ?",
                view_count_int + 1, other_set["doc_name"], now_date,
            )
        }
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    return return_data
}
