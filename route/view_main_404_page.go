package route

import (
	"opennamu/route/tool"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
)

func View_main_404_page(config tool.Config) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    page_404_set := ""
    tool.QueryRow_DB(
        db,
        `select data from other where name = "manage_404_page"`,
        []any{ &page_404_set },
    )

    return_data := make(map[string]any)
    return_data["response"] = "ok"

    page_404_dir := filepath.Join("..", "404.html")
    if tool.File_exist_check(page_404_dir) && page_404_set == "404_file" {
        return_data["data"] = tool.File_text_read(page_404_dir)
    } else {
        db_data := ""

        tool.QueryRow_DB(
            db,
            `select data from other where name = "manage_404_page_content"`,
            []any{ &db_data },
        )

        if db_data != "" {
            return_data["data"] = tool.Get_template(db, config, "404", db_data)
        } else {
            return_data["data"] = tool.Get_template(db, config, "404", tool.Get_language(db, "func_404_error", true))
        }
    }

    json_data, _ := json.Marshal(return_data)

    data := tool.View_result{
        HTML : return_data["data"].(string),
        JSON : string(json_data),
    }

    return data
}