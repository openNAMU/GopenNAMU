package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_bbs_in(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    bbs_name := ""

    tool.QueryRow_DB(
        db,
        tool.DB_change("select set_data from bbs_set where set_id = ? and set_name = 'bbs_name'"),
        []any{ &bbs_name },
        other_set["bbs_num"],
    )

    out := tool.Get_template(
        db,
        config,
        bbs_name,
        "",
        "",
        [][]any{},
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = out

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}