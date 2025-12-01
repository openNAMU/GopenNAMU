package tool

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

var db_set = map[string]string{}

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

func DB_init_standalone() {
    db_env, db_env_exist := os.LookupEnv("NAMU_DB")
    db_env_type, db_env_type_exist := os.LookupEnv("NAMU_DB_TYPE")
    if db_env_exist || db_env_type_exist {
        db_set["db_name"] = Choose(db_env, "data")
        db_set["db_type"] = Choose(db_env_type, "sqlite")

        return
    }
    
    path_dir := filepath.Join("..", "data", "set.json")
    if File_exist_check(path_dir) {
        raw, err := os.ReadFile(path_dir)
        if err == nil {
            tmp := map[string]string{}
            if err := json.Unmarshal(raw, &tmp); err == nil {
                if v, ok := tmp["db_name"]; ok {
                    db_set["db_name"] = v
                } else {
                    db_set["db_name"] = "data"
                }
                
                if v, ok := tmp["db_type"]; ok {
                    db_set["db_type"] = v
                } else {
                    db_set["db_type"] = "sqlite"
                }

                return
            }
        }
    }

	db_set["db_name"] = "data"
	db_set["db_type"] = "sqlite"
    
    return
}

func DB_init_standalone_MySQL() {
    path := filepath.Join("..", "data", "mysql.json")
	if !File_exist_check(path) {
		return
	}

    raw, err := os.ReadFile(path)
	if err != nil {
		return
	}

	tmp := map[string]string{}
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return
	}

	if host, ok := tmp["host"]; ok && host != "" {
		db_set["db_mysql_host"] = host
	} else {
		db_set["db_mysql_host"] = "localhost"
	}

	if port, ok := tmp["port"]; ok && port != "" {
		db_set["db_mysql_port"] = port
	} else {
		db_set["db_mysql_port"] = "3306"
	}

	if user, ok := tmp["user"]; ok {
		db_set["db_mysql_user"] = user
	}

	if pw, ok := tmp["password"]; ok {
		db_set["db_mysql_pw"] = pw
	}

    return
}

func Exec_DB(db *sql.DB, query string, values ...any) {
    const retryDelay = 10 * time.Millisecond

    stmt, err := db.Prepare(DB_change(query))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    for {
        _, err = stmt.Exec(values...)
        if err == nil {
            return
        }

        if strings.Contains(err.Error(), "database is locked") {
            time.Sleep(retryDelay)
            continue
        }

        panic(err)
    }
}

func Query_DB(db *sql.DB, query string, values ...any) *sql.Rows {
    const retryDelay = 10 * time.Millisecond

    stmt, err := db.Prepare(DB_change(query))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    for {
        rows, err := stmt.Query(values...)
        if err == nil {
            return rows
        }

        if strings.Contains(err.Error(), "database is locked") {
            time.Sleep(retryDelay)
            continue
        }

        panic(err)
    }
}

// 이래서 포인터를 배우는구나...
func QueryRow_DB(db *sql.DB, query string, var_list []any, values ...any) bool {
    const retryDelay = 10 * time.Millisecond

    stmt, err := db.Prepare(DB_change(query))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    for {
        row := stmt.QueryRow(values...)

        err := row.Scan(var_list...)
        switch err {
        case nil:
            return true
        case sql.ErrNoRows:
            return false
        }

        if strings.Contains(err.Error(), "database is locked") {
            time.Sleep(retryDelay)
            continue
        }

        panic(err)
    }
}

func DB_init() {
    DB_init_standalone()

    if db_set["db_type"] == "mysql" {
        DB_init_standalone_MySQL()
    }
}

func DB_connect() *sql.DB {
    // log.Default().Println("DB open")

    if db_set["db_type"] == "sqlite" {
        db, err := sql.Open("sqlite", filepath.Join("..", db_set["db_name"] + ".db") + "?_journal_mode=WAL&_busy_timeout=5000")
        if err != nil {
            panic(err)
        }

        return db
    } else {
        db, err := sql.Open("mysql", db_set["db_mysql_user"] + ":" + db_set["db_mysql_pw"] + "@tcp(" + db_set["db_mysql_host"] + ":" + db_set["db_mysql_port"] + ")/" + db_set["db_name"])
        if err != nil {
            panic(err)
        }

        return db
    }
}

func DB_close(db *sql.DB) {
    db.Close()
    
    // log.Default().Println("DB close")
}

func Get_DB_type() string {
    return db_set["db_type"]
}

func DB_change(data string) string {
    if Get_DB_type() == "mysql" {
        data = strings.Replace(data, "random()", "rand()", -1)
        data = strings.Replace(data, "collate nocase", "collate utf8mb4_general_ci", -1)
    }

    return data
}
