package tool

import (
	"database/sql"
	"strconv"
)

func Send_alarm(db *sql.DB, from string, target string, data string) {
    if from != target {
        data = from + " | " + data

        now_time := Get_time()

        var count string

        stmt, err := db.Prepare(DB_change("select id from user_notice where name = ? order by id + 0 desc limit 1"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        err = stmt.QueryRow(target).Scan(&count)
        if err != nil {
            if err == sql.ErrNoRows {
                count = "1"
            } else {
                panic(err)
            }
        }

        count_int, _ := strconv.Atoi(count)
        count_int += 1

        Exec_DB(
            db,
            "insert into user_notice (id, name, data, date, readme) values (?, ?, ?, ?, '')",
            count_int, target, data, now_time,
        )
    }
}
