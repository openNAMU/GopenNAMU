package tool

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	_ "modernc.org/sqlite"
)

var db_set = map[string]string{}

func Exec_DB(db *sql.DB, query string, values ...interface{}) {
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

func DB_init(get_db_set string) {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(get_db_set), &other_set)

    for k, v := range other_set {
        db_set[k] = v
    }
}

func DB_connect() *sql.DB {
    // log.Default().Println("DB open")

    if db_set["db_type"] == "sqlite" {
        db, err := sql.Open("sqlite", db_set["db_name"] + ".db?_journal_mode=WAL&_busy_timeout=5000")
        if err != nil {
            panic(err)
        }

        /*
        var journal_mode string

        err = db.QueryRow("PRAGMA journal_mode").Scan(&journal_mode)
        if err != nil {
            panic(err)
        }

        log.Default().Println(journal_mode)
        */

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
