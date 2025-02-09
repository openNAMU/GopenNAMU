package tool

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html"
	"html/template"
	"net/url"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
)

func Sha224(data string) string {
    hasher := sha256.New224()
    hasher.Write([]byte(data))
    hash_byte := hasher.Sum(nil)
    hash_str := hex.EncodeToString(hash_byte)

    return hash_str
}

func Url_parser(data string) string {
    return url.QueryEscape(data)
}

func HTML_escape(data string) string {
    return template.HTMLEscapeString(data)
}

func HTML_unescape(data string) string {
    return html.UnescapeString(data)
}

func Arr_in_str(arr []string, data string) bool {
    for _, v := range arr {
        if v == data {
            return true
        }
    }

    return false
}

func Get_time() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

func Get_date() string {
    return time.Now().Format("2006-01-02")
}

func Get_month() string {
    return time.Now().Format("2006-01")
}

func Get_document_setting(db *sql.DB, doc_name string, set_name string, doc_rev string) [][]string {
    var rows *sql.Rows

    if doc_rev != "" {
        stmt, err := db.Prepare(DB_change("select set_data, doc_rev from data_set where doc_name = ? and doc_rev = ? and set_name = ?"))
        if err != nil {
            panic(err)
        }

        defer stmt.Close()

        rows, err = stmt.Query(doc_name, doc_rev, set_name)
        if err != nil {
            panic(err)
        }
    } else {
        stmt, err := db.Prepare(DB_change("select set_data, doc_rev from data_set where doc_name = ? and set_name = ?"))
        if err != nil {
            panic(err)
        }

        defer stmt.Close()

        rows, err = stmt.Query(doc_name, set_name)
        if err != nil {
            panic(err)
        }
    }
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var set_data string
        var doc_rev string

        err := rows.Scan(&set_data, &doc_rev)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{set_data, doc_rev})
    }

    return data_list
}

func Get_setting(db *sql.DB, set_name string, data_coverage string) [][]string {
    var rows *sql.Rows

    if data_coverage != "" {
        stmt, err := db.Prepare(DB_change("select data, coverage from other where name = ? and coverage = ?"))
        if err != nil {
            panic(err)
        }

        defer stmt.Close()

        rows, err = stmt.Query(set_name, data_coverage)
        if err != nil {
            panic(err)
        }
    } else {
        stmt, err := db.Prepare(DB_change("select data, coverage from other where name = ?"))
        if err != nil {
            panic(err)
        }

        defer stmt.Close()

        rows, err = stmt.Query(set_name)
        if err != nil {
            panic(err)
        }
    }
    defer rows.Close()

    data_list := [][]string{}

    for rows.Next() {
        var set_data string
        var set_coverage string

        err := rows.Scan(&set_data, &set_coverage)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, []string{set_data, set_coverage})
    }

    return data_list
}

func Get_skin_list(data string, default_flag bool) []string {
    entries, err := os.ReadDir("views")
	if err != nil {
		return nil
	}

	var skin_list []string

	if default_flag {
		skin_list = append(skin_list, "default")
	}

	for _, entry := range entries {
		skin_list = append(skin_list, entry.Name())
	}

	var skin_return_data []string

	for _, skin_data := range skin_list {
		if skin_data != "main_css" {
			if skin_data == data {
				skin_return_data = append([]string{skin_data}, skin_return_data...)
			} else {
				skin_return_data = append(skin_return_data, skin_data)
			}
		}
	}

	return skin_return_data
}

func Get_use_skin_name(db *sql.DB, ip string) string {
    skin_list := Get_skin_list("ringo", true)
    skin := skin_list[0]

    user_skin_name := ""
    if IP_or_user(ip) {
        stmt, err := db.Prepare(DB_change("select data from user_set where name = 'skin' and id = ?"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        err = stmt.QueryRow(ip).Scan(&user_skin_name)
        if err != nil {
            if err == sql.ErrNoRows {
            } else {
                panic(err)
            }
        }
    }

    if user_skin_name == "default" {
        user_skin_name = ""
    }

    if user_skin_name == "" {
        stmt, err := db.Prepare(DB_change("select data from other where name = 'skin'"))
        if err != nil {
            panic(err)
        }
        defer stmt.Close()

        err = stmt.QueryRow().Scan(&user_skin_name)
        if err != nil {
            if err == sql.ErrNoRows {
            } else {
                panic(err)
            }
        }
    }

    if user_skin_name != "" && Arr_in_str(skin_list, user_skin_name) {
        skin = user_skin_name
    }

    return skin
}

func Get_template(db *sql.DB, ip string, data jet.VarMap) string {
    views := jet.NewSet(
		jet.NewOSFileSystemLoader("./views/" + Get_use_skin_name(db, ip)),
		jet.InDevelopmentMode(),
	)

    tmpl, err := views.GetTemplate("example.jet")
	if err != nil {
		panic(err)
	}

    var buf bytes.Buffer

    err = tmpl.Execute(&buf, data, nil)
	if err != nil {
        panic(err)
    }

    return buf.String()
}

func Get_wiki_set(db *sql.DB, ip string) {
    skin_name := Get_use_skin_name(db, ip)
    data_list := []string{}
    
    set_wiki_name := ""
    set_license := ""
    set_logo := ""

    stmt, err := db.Prepare(DB_change("select data from other where name = 'name'"))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    err = stmt.QueryRow().Scan(&set_wiki_name)
    if err != nil {
        if err == sql.ErrNoRows {
            set_wiki_name = "Wiki"
        } else {
            panic(err)
        }
    }

    stmt2, err := db.Prepare(DB_change("select data from other where name = 'license'"))
    if err != nil {
        panic(err)
    }
    defer stmt2.Close()

    err = stmt2.QueryRow().Scan(&set_license)
    if err != nil {
        if err == sql.ErrNoRows {
        } else {
            panic(err)
        }
    }

    stmt3, err := db.Prepare(DB_change("select data from other where name = 'logo' and coverage = ?"))
    if err != nil {
        panic(err)
    }
    defer stmt3.Close()

    err = stmt3.QueryRow(skin_name).Scan(&set_logo)
    if err != nil {
        if err == sql.ErrNoRows {
        } else {
            panic(err)
        }
    }

    if set_logo == "" {
        stmt4, err := db.Prepare(DB_change("select data from other where name = 'logo' and coverage = ''"))
        if err != nil {
            panic(err)
        }
        defer stmt4.Close()

        err = stmt4.QueryRow().Scan(&set_logo)
        if err != nil {
            if err == sql.ErrNoRows {
            } else {
                panic(err)
            }
        }
    }
    
    data_list = append(data_list, set_wiki_name)
    data_list = append(data_list, set_license)
    data_list = append(data_list, set_logo)
}