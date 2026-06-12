package tool

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func DB_table_list() map[string][]string {
    create_data := map[string][]string{}

    // 폐지 예정 (data_set으로 통합)
    create_data["data_set"] = []string{"doc_name", "doc_rev", "set_name", "set_data"}

    create_data["data"] = []string{"title", "data", "type"}
    create_data["history"] = []string{"id", "title", "data", "date", "ip", "send", "leng", "hide", "type"}
    create_data["rc"] = []string{"id", "title", "date", "type"}
    create_data["acl"] = []string{"title", "data", "type"}

    // 개편 예정 (data_link로 변경)
    create_data["back"] = []string{"title", "link", "type", "data"}

    // 폐지 예정 (topic_set으로 통합) [가장 시급]
    create_data["topic_set"] = []string{"thread_code", "set_name", "set_id", "set_data"}

    create_data["rd"] = []string{"title", "sub", "code", "date", "band", "stop", "agree", "acl"}
    create_data["topic"] = []string{"id", "data", "date", "ip", "block", "top", "code"}

    // 폐지 예정 (user_set으로 통합)
    create_data["rb"] = []string{"block", "end", "today", "blocker", "why", "band", "login", "ongoing"}

    // 개편 예정 (wiki_set과 wiki_filter과 wiki_vote으로 변경)
    create_data["other"] = []string{"name", "data", "coverage"}
    create_data["html_filter"] = []string{"html", "kind", "plus", "plus_t"}
    create_data["vote"] = []string{"name", "id", "subject", "data", "user", "type", "acl"}

    // 개편 예정 (auth와 auth_log로 변경)
    create_data["alist"] = []string{"name", "acl"}
    create_data["re_admin"] = []string{"who", "what", "time"}

    // 개편 예정 (user_notice와 user_agent로 변경)
    create_data["ua_d"] = []string{"name", "ip", "ua", "today", "sub"}

    create_data["user_set"] = []string{"name", "id", "data"}
    create_data["user_notice"] = []string{"id", "name", "data", "date", "readme"}

    create_data["bbs_set"] = []string{"set_name", "set_code", "set_id", "set_data"}
    create_data["bbs_data"] = []string{"set_name", "set_code", "set_id", "set_data"}

    return create_data
}

func DB_make_MySQL(db *sql.DB, new_db_set map[string]string) {
    Exec_DB(
        db,
        `create database if not exists ` + new_db_set["db_name"] + ` default character set utf8mb4`,
    )
}

func DB_column_exists(db *sql.DB, table_name string, column_name string) bool {
    query := fmt.Sprintf(
        "SELECT %s FROM %s LIMIT 1",
        column_name,
        table_name,
    )

    rows, err := db.Query(query)
    if err != nil {
        return false
    }
    defer rows.Close()

    return true
}

func DB_field_text(db_type string) string {
    if db_type == "mysql" {
        return "longtext"
    }

    return "text default ''"
}


func DB_warn_null_column(db *sql.DB, table_name string, column_name string) {
    query := fmt.Sprintf(
        "select count(*) from %s where %s is null",
        table_name,
        column_name,
    )

    var count int
    err := db.QueryRow(query).Scan(&count)
    if err != nil {
        log.Printf("[DB WARNING] null check failed: %s.%s: %v", table_name, column_name, err)
        return
    }

    if count > 0 {
        log.Printf("[DB WARNING] %s.%s has %d null values", table_name, column_name, count)
    }
}


func DB_create_table(db *sql.DB, table_name string, field_text string) {
    Exec_DB(
        db,
        fmt.Sprintf(
            "create table if not exists %s (test %s)",
            table_name,
            field_text,
        ),
    )
}

func DB_alter_add_column(db *sql.DB, table_name string, column_name string, field_text string) {
    Exec_DB(
        db,
        fmt.Sprintf(
            "alter table %s add column %s %s",
            table_name,
            column_name,
            field_text,
        ),
    )
}

func DB_make(db *sql.DB, new_db_set map[string]string) error {
    if new_db_set["db_type"] == "mysql" {
        DB_make_MySQL(db, new_db_set)
    } else {
        Exec_DB(
            db,
            `pragma journal_mode = WAL`,
        )
    }

    db = DB_connect()
    defer DB_close(db)

    if err := db.Ping(); err != nil {
        return err
    }

    field_text := DB_field_text(new_db_set["db_type"])

    table_list := DB_table_list()
    for table_name, table_data := range table_list {
        DB_create_table(db, table_name, field_text)

        columns := append([]string{"test"}, table_data...)

        for _, column_name := range columns {
            if !DB_column_exists(db, table_name, column_name) {
                DB_alter_add_column(db, table_name, column_name, field_text)
            }

            DB_warn_null_column(db, table_name, column_name)
        }
    }

    return nil
}

func DB_init() {
    new_db_set := DB_boot()
    
    db, err := DB_connect_init()
    if err != nil {
        panic(fmt.Errorf("DB connection failed: %w", err))
    }
    defer DB_close(db)

    if err := DB_make(db, new_db_set); err != nil {
        panic(fmt.Errorf("DB setup failed: %w", err))
    }
}

func Main_init() {
    DB_init()
    DB_boot()

    db := DB_connect()
    defer DB_close(db)

    now_version := ""
    QueryRow_DB(
        db,
        `select data from other where name = "ver"`,
        []any{ &now_version },
    )

    if now_version == "" {
        First_init()
    } else {
        last_version := Get_last_version()

        if now_version != last_version["c_ver"] {
            Make_set_json()
        }
    }
}

func Get_last_version() map[string]string {
    version_file_path := filepath.Join("..", "version.json")
    if _, err := os.Stat(version_file_path); err == nil {
        data, err := os.ReadFile(version_file_path)
        if err != nil {
            panic(err)
        }

        version_json := map[string]string{}
        json.Unmarshal([]byte(data), &version_json)

        return version_json
    } else {
        panic(err)
    }
}

func First_init() {

}

func Make_set_json() {

}