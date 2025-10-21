package route

import (
	"opennamu/route/tool"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Api_edit_post_exter(config tool.Config) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_edit_post(
        config,
        other_set["doc_name"],
        other_set["data"],
        other_set["send"],
        other_set["agree"],
    )

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_edit_post(config tool.Config, doc_name string, data string, send string, agree string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)
    
    date := tool.Get_date()
    data = strings.ReplaceAll(data, "\r", "")

    if !tool.Do_edit_slow_check(db, config, "edit") {
        return_data["response"] = "error"
        return_data["data"] = "slow edit limit"

        return return_data
    } else if !tool.Do_edit_filter(db, config, doc_name, data) {
        return_data["response"] = "error"
        return_data["data"] = "edit filter (content)"

        return return_data
    } else if !tool.Do_edit_filter(db, config, doc_name, send) {
        return_data["response"] = "error"
        return_data["data"] = "edit filter (send)"

        return return_data
    } else if !tool.Do_edit_send_require_check(db, config, send) {
        return_data["response"] = "error"
        return_data["data"] = "send require"

        return return_data
    } else if !tool.Do_edit_text_checkbox_check(db, config, agree) {
        return_data["response"] = "error"
        return_data["data"] = "checkbox check require"

        return return_data
    } else if !tool.Do_edit_max_length_check(db, config, data) {
        return_data["response"] = "error"
        return_data["data"] = "overflow max length"

        return return_data
    }

    var old_data string

    tool.Query_DB(
        db,
        `select data from data where title = ?`,
        []any{ doc_name },
    )

    length := tool.Get_edit_length_diff(old_data, data)

    tool.Do_add_history(
        db,
        doc_name,
        data,
        date,
        config.IP,
        send,
        length,
        "",
        "",
    )

    return_data["response"] = "ok"

    return return_data
}