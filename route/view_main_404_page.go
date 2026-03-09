package route

import (
	"opennamu/route/tool"
	"path/filepath"
)

func View_main_404_page(config tool.Config, url string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if url == "/" {
        frontpage := "FrontPage"

        tool.QueryRow_DB(
            db,
            `select data from other where name = "frontpage"`,
            []any{ &frontpage },
        )

        return tool.Get_redirect("/w/" + tool.Url_parser(frontpage))
    }

    page_404_set := ""
    tool.QueryRow_DB(
        db,
        `select data from other where name = "manage_404_page"`,
        []any{ &page_404_set },
    )

    data_html := ""

    page_404_dir := filepath.Join("..", "404.html")
    if tool.File_exist_check(page_404_dir) && page_404_set == "404_file" {
        data_html = tool.File_text_read(page_404_dir)
    } else {
        db_data := ""

        tool.QueryRow_DB(
            db,
            `select data from other where name = "manage_404_page_content"`,
            []any{ &db_data },
        )

        if db_data != "" {
            data_html = tool.Get_template(
                db,
                config,
                "404",
                db_data,
                []any{},
                [][]any{},
            )
        } else {
            data_html = tool.Get_template(
                db,
                config,
                "404",
                tool.Get_language(db, "func_404_error", true),
                []any{},
                [][]any{},
            )
        }
    }

    return data_html
}