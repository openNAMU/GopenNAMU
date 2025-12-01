package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_bbs_w_comment_make(db *sql.DB, doc_name string) string {
    inter_other_set := map[string]string{}
    inter_other_set["set_id"] = "0"
    inter_other_set["title"] = doc_name
    inter_other_set["data"] = ""

    json_data, _ := json.Marshal(inter_other_set)

    send_request := tool.Config{
        Other_set: string(json_data),
        IP: "tool:system",
    }

    return_data := Api_bbs_w_comment_one(send_request, false)

    return_data_api := map[string]string{}
    json.Unmarshal([]byte(return_data), &return_data_api)

    return return_data_api["data"]
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