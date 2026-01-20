package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_bbs_w_comment_make(db *sql.DB, doc_name string) string {
    return ""
}

func Api_w_comment(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    db_code := tool.Get_document_setting(db, other_set["doc_name"], "document_comment_code", "")
    
    db_code_str := ""
    if len(db_code) >= 1 {
        db_code_str = db_code[0][0]
    }

    if db_code_str == "" {
        db_code_str = Api_bbs_w_comment_make(db, other_set["doc_name"])

        tool.Exec_DB(
            db,
            "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_comment_code', ?)",
            other_set["doc_name"], db_code_str,
        )
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = db_code_str

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}