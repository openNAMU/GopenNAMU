package tool

import (
	"database/sql"

	"github.com/dlclark/regexp2"
)

func Do_edit_filter(db *sql.DB, config Config, doc_name string, data string) bool {
    if !Check_acl(db, "", "", "edit_filter_pass", config.IP) {
        rows := Query_DB(
            db, 
            `select plus, plus_t from html_filter where kind = 'regex_filter' and plus != ''`,
        )
        defer rows.Close()

        for rows.Next() {
            var plus string
            var plus_t string

            err := rows.Scan(&plus, &plus_t)
            if err != nil {
                panic(err)
            }

            r, err := regexp2.Compile(plus, 0)
            if err != nil {
                continue
            }

            m, err := r.MatchString(data)
            if err == nil && m {
                return false
            }
        }
    }

    return true
}