package route

import (
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_post(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
        return_data := make(map[string]any)
        return_data["response"] = "require auth"

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }

    set_code := ""
    tool.QueryRow_DB(
        db,
        tool.DB_change("select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc"),
        []any{ &set_code },
        other_set["set_id"],
    )

    set_code_int, _ := strconv.Atoi(set_code)
    set_code_int += 1

    set_code_str := strconv.Itoa(set_code_int)

    date_now := tool.Get_time()

    insert_db := [][]string{
        { "title", other_set["title"] },
        { "data", other_set["data"] },
        { "date", date_now },
        { "user_id", config.IP },
    }
    for _, v := range insert_db {
        tool.Exec_DB(
            db,
            "insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
            v[0], set_code_str, other_set["set_id"], v[1],
        )
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = set_code_str

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}