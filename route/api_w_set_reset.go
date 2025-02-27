package route

import (
	"database/sql"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_set_reset(db *sql.DB, config tool.Config) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set[0]), &other_set)

    doc_name := other_set["name"]
    ip := config.IP

    if tool.Check_acl(db, "", "", "owner_auth", ip) {
        tool.Exec_DB(
            db,
            "delete from acl where title = ?",
            doc_name,
        )

        tool.Exec_DB(
            db,
            "delete from data_set where doc_name = ? and set_name = 'acl_date'",
            doc_name,
        )

        set_list := []string{
            "document_markup",
            "document_top",
            "document_editor_top",
        }

        for for_a := 0; for_a < len(set_list); for_a++ {
            tool.Exec_DB(
                db,
                "delete from data_set where doc_name = ? and set_name = ?",
                doc_name, set_list[for_a],
            )
        }

        return_data := make(map[string]interface{})
        return_data["response"] = "ok"
        return_data["language"] = map[string]string{
            "reset": tool.Get_language(db, "reset", false),
        }

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    } else {
        return_data := make(map[string]interface{})
        return_data["response"] = "require auth"
        return_data["language"] = map[string]string{
            "authority_error": tool.Get_language(db, "authority_error", false),
        }

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
}
