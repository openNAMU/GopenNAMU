package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var standalone_mode = false

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

func Get_IP(c *gin.Context) string {
    if standalone_mode {
        return c.ClientIP()
    } else {
        return c.Request.Header.Get("X-Forwarded-For")
    }
}

func Get_Cookies(c *gin.Context) string {
    return c.Request.Header.Get("Cookie")
}

func Get_session(c *gin.Context) string {
    return ""
}

func Get_document_setting(db *sql.DB, doc_name string, set_name string, doc_rev string) [][]string {
    var rows *sql.Rows
    if doc_rev != "" {
        rows = Query_DB(
            db,
            DB_change("select set_data, doc_rev from data_set where doc_name = ? and doc_rev = ? and set_name = ?"),
            doc_name, doc_rev, set_name,
        )
    } else {
        rows = Query_DB(
            db,
            DB_change("select set_data, doc_rev from data_set where doc_name = ? and set_name = ?"),
            doc_name, set_name,
        )
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
        rows = Query_DB(
            db,
            DB_change("select data, coverage from other where name = ? and coverage = ?"),
            set_name, data_coverage,
        )
    } else {
        rows = Query_DB(
            db,
            DB_change("select data, coverage from other where name = ?"),
            set_name,
        )
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
    entries, err := os.ReadDir(filepath.Join("..", "views"))
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
        QueryRow_DB(
            db,
            DB_change("select data from user_set where name = 'skin' and id = ?"),
            []any{ &user_skin_name },
            ip,
        )
    }

    if user_skin_name == "default" {
        user_skin_name = ""
    }

    if user_skin_name == "" {
        QueryRow_DB(
            db,
            DB_change("select data from other where name = 'skin'"),
            []any{ &user_skin_name },
        )
    }

    if user_skin_name != "" && Arr_in_str(skin_list, user_skin_name) {
        skin = user_skin_name
    }

    return skin
}

func Get_skin_route(db *sql.DB, ip string) string {
    return filepath.Join("..", "views", Get_use_skin_name(db, ip), "index.html")
}

func Get_domain(db *sql.DB, full_string bool) string {
    domain := ""
    sys_host := ""

    if full_string {
        http_select := ""
        QueryRow_DB(
            db,
            DB_change("select data from other where name = 'http_select'"),
            []any{ &http_select },
        )
        
        if http_select == "" {
            http_select = "http"
        }

        domain = http_select + "://"

        db_domain := ""
        QueryRow_DB(
            db,
            DB_change("select data from other where name = 'domain'"),
            []any{ &db_domain },
        )

        if db_domain != "" {
            domain += db_domain
        } else {
            domain += sys_host
        }
    } else {
        db_domain := ""
        QueryRow_DB(
            db,
            DB_change("select data from other where name = 'domain'"),
            []any{ &db_domain },
        )

        if db_domain != "" {
            domain = db_domain
        } else {
            domain = sys_host
        }
    }

    return domain
}

func Get_wiki_custom(db *sql.DB, ip string, session_str string, cookies string) []any {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    
    session := map[string]string{}
    json.Unmarshal([]byte(session_str), &session)

    skin_name := "_" + Get_use_skin_name(db, ip)

    user_icon := 1
    user_name := ip
    user_head := ""
    user_email := ""
    user_admin := "0"
    user_acl_list := []string{}
    user_notice_count := 0

    if !IP_or_user(ip) {
        user_head_main := ""
        QueryRow_DB(
            db,
            DB_change("select data from user_set where id = ? and name = 'custom_css'"),
            []any{ &user_head_main },
            ip,
        )

        user_head_skin := ""
        QueryRow_DB(
            db,
            DB_change("select data from user_set where id = ? and name = ?"),
            []any{ &user_head_main },
            ip,
            "custom_css" + skin_name,
        )

        user_head += user_head_main + user_head_skin

        QueryRow_DB(
            db,
            DB_change("select data from user_set where name = 'email' and id = ?"),
            []any{ &user_email },
            ip,
        )

        if Check_acl(db, "", "", "all_admin_auth", ip) {
            user_admin = "1"

            acl_name := ""
            QueryRow_DB(
                db,
                DB_change("select data from user_set where id = ? and name = 'acl'"),
                []any{ &acl_name },
                ip,
            )

            rows := Query_DB(
                db,
                DB_change("select acl from alist where name = ?"),
                acl_name,
            )
            defer rows.Close()

            for rows.Next() {
                user_acl_name := ""

                err := rows.Scan(&user_acl_name)
                if err != nil {
                    panic(err)
                }

                user_acl_list = append(user_acl_list, user_acl_name)
            }
        }

        QueryRow_DB(
            db,
            DB_change("select count(*) from user_notice where name = ? and readme = ''"),
            []any{ &user_notice_count },
            ip,
        )
    } else {
        user_icon = 0
        user_name = Get_language(db, "user", true)
        user_email = ""
        user_acl_list = []string{}
        user_notice_count = 0
        user_head = ""
    }

    user_ban := "0"
    user_ban_check := Get_user_ban(db, ip, "")[0]
    if user_ban_check == "true" {
        user_ban = "1"
    }

    user_topic := "0"
    user_topic_check := QueryRow_DB(
        db,
        DB_change("select title from rd where title = ? and stop = '' limit 1"),
        []any{},
        "user:" + ip,
    )
    if user_topic_check {
        user_topic = "1"
    }

    return []any{
        "",
        "",
        user_icon,
        user_head,
        user_email,
        user_name,
        user_admin,
        user_ban,
        user_notice_count,
        func(user_acl_list []string) any {
            if len(user_acl_list) == 0 {
                return "0"
            } else {
                return user_acl_list
            }
        }(user_acl_list),
        ip,
        user_topic,
        "",
        Get_level(db, ip),
    }
}

func Get_wiki_set(db *sql.DB, ip string, cookies string) []any {
    skin_name := Get_use_skin_name(db, ip)
    data_list := []any{}

    set_wiki_name := "Wiki"
    QueryRow_DB(
        db,
        DB_change("select data from other where name = 'name'"),
        []any{ &set_wiki_name },
    )

    set_license := ""
    QueryRow_DB(
        db,
        DB_change("select data from other where name = 'license'"),
        []any{ &set_license },
    )

    set_logo := ""
    QueryRow_DB(
        db,
        DB_change("select data from other where name = 'logo' and coverage = ?"),
        []any{ &set_logo },
        skin_name,
    )

    if set_logo == "" {
        QueryRow_DB(
            db,
            DB_change("select data from other where name = 'logo' and coverage = ''"),
            []any{ &set_logo },
        )
    }

    if set_logo == "" {
        set_logo = set_wiki_name
    }
    
    set_head := ""
    QueryRow_DB(
        db,
        DB_change("select data from other where name = 'head' and coverage = ''"),
        []any{ &set_head },
    )

    set_head_skin := ""
    QueryRow_DB(
        db,
        DB_change("select data from other where name = 'head' and coverage = ?"),
        []any{ &set_head_skin },
        skin_name,
    )

    set_head_dark := ""

    cookie_map := Get_cookie_header(cookies)
    if cookie_map["main_css_darkmode"] == "1" {
        QueryRow_DB(
            db,
            DB_change("select data from other where name = 'head' and coverage = ?"),
            []any{ &set_head_dark },
            skin_name + "-cssdark",
        )
    }
    
    set_top_menu := ""
    QueryRow_DB(
        db,
        DB_change("select data from other where name = 'top_menu'"),
        []any{ &set_top_menu },
    )

    set_top_menu_user := ""
    QueryRow_DB(
        db,
        DB_change("select data from user_set where name = 'top_menu' and id = ?"),
        []any{ &set_top_menu_user },
        ip,
    )

    set_top_menu = strings.ReplaceAll(set_top_menu, "\r", "")
    set_top_menu_user = strings.ReplaceAll(set_top_menu_user, "\r", "")

    set_top_menu_mix := ""
    if set_top_menu != "" && set_top_menu_user != "" {
        set_top_menu_mix = set_top_menu + "\n" + set_top_menu_user
    } else {
        set_top_menu_mix = set_top_menu + set_top_menu_user
    }

    set_top_menu_result := [][]string{}
    if set_top_menu_mix != "" {
        lst := strings.Split(set_top_menu_mix, "\n")
        if len(lst) % 2 != 0 {
            lst = append(lst, "")
        }

        for i := 0; i < len(lst) - 1; i += 2 {
            set_top_menu_result = append(set_top_menu_result, []string{lst[i], lst[i+1]})
        }
    }

    template_var := []any{}
    for for_a := 1; for_a < 4; for_a++ {
        template_var_tmp := ""
        QueryRow_DB(
            db,
            DB_change("select data from other where name = ?"),
            []any{ &template_var_tmp },
            "template_var_" + strconv.Itoa(for_a),
        )

        template_var = append(template_var, template_var_tmp)
    }
    
    data_list = append(data_list, set_wiki_name)
    data_list = append(data_list, set_license)
    data_list = append(data_list, "")
    data_list = append(data_list, "")
    data_list = append(data_list, set_logo)
    data_list = append(data_list, set_head + set_head_skin + set_head_dark)
    
    if len(set_top_menu_result) > 0 {
        data_list = append(data_list, set_top_menu_result)
    } else {
        data_list = append(data_list, "")
    }

    data_list = append(data_list, template_var...)

    return data_list
}

func Get_cookie_header(cookie_header string) map[string]string {
    cookies := make(map[string]string)
    
    parts := strings.Split(cookie_header, ";")
    for _, part := range parts {
        part = strings.TrimSpace(part)
        if len(part) == 0 {
            continue
        }

        kv := strings.SplitN(part, "=", 2)
        if len(kv) == 2 {
            key := strings.TrimSpace(kv[0])
            value := strings.TrimSpace(kv[1])
            
            cookies[key] = value
        }
    }

    return cookies
}

type Config struct {
    Other_set string
    IP string
    Cookies string
    Session string
}

func Deep_copy_config(src Config) Config {
    return Config{
        Other_set : strings.Clone(src.Other_set),
        IP : strings.Clone(src.IP),
        Cookies : strings.Clone(src.Cookies),
        Session : strings.Clone(src.Session),
    }
}

func Choose(v string, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}

	return v
}

func File_exist_check(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func File_text_read(path string) string {
    raw, err := os.ReadFile(path)
    if err != nil {
        return ""
    }

    return string(raw)
}

func IN_mod_OUT_mod(data bool) {
    standalone_mode = data
}

func Str_to_int(data string) int {
    num, _ := strconv.Atoi(data)
    return num
}