package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html"
	"html/template"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
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

func Get_IP(c *gin.Context) string {
    return c.Request.Header.Get("X-Forwarded-For")
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
    return "./views/" + Get_use_skin_name(db, ip) + "/index.html"
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

func Cache_v() string {
    return ".cache_v288"
}

func Get_wiki_css(data []any, cookies string) []any {
    for len(data) < 4 {
        data = append(data, "")
    }

    data_css := ""
    data_css_dark := ""

    data_css_ver := Cache_v()

    // Cache Control
    data_css += `<meta http-equiv="Cache-Control" content="max-age=31536000">`

    // External JS
    data_css += `<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js" integrity="sha384-7zkQWkzuo3B5mTepMUcHkMB5jZaolc2xDwL6VFqjFALcbeS9Ggm/Yr2r3Dy4lfFg" crossorigin="anonymous"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js" integrity="sha512-rdhY3cbXURo13l/WU9VlaRyaIYeJ/KBakckXIvJNAQde8DgpOmE+eZf7ha4vdqVjTtwQt69bD2wH2LXob/LB7Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/languages/x86asm.min.js" integrity="sha512-HeAchnWb+wLjUb2njWKqEXNTDlcd1QcyOVxb+Mc9X0bWY0U5yNHiY5hTRUt/0twG8NEZn60P3jttqBvla/i2gA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.48.0/min/vs/loader.min.js" integrity="sha512-ZG31AN9z/CQD1YDDAK4RUAvogwbJHv6bHrumrnMLzdCrVu4HeAqrUX7Jsal/cbUwXGfaMUNmQU04tQ8XXl5Znw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/highlightjs-line-numbers.js/2.8.0/highlightjs-line-numbers.min.js"></script>`

    // Func JS
    data_css += `<script defer src="/views/main_css/js/func/func.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_version.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_user_info.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_version_skin.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_http_warning_text.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/ie_end_of_life.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/shortcut.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/editor.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/render.js` + data_css_ver + `"></script>`

    // Main CSS
    data_css += `<link rel="stylesheet" href="/views/main_css/css/main.css` + data_css_ver + `">`

    // External CSS
    data_css += `<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css" integrity="sha384-nB0miv6/jRmo5UMMR1wu3Gz6NLsoTkbqJghGIsx//Rlm+ZU03BU6SQNC66uf4l5+" crossorigin="anonymous">`
    data_css += `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/default.min.css" integrity="sha512-hasIneQUHlh06VNBe7f6ZcHmeRTLIaQWFd43YriJ0UND19bvYRauxthDg8E4eVNPm9bRUhr5JGeqH7FRFXQu5g==" crossorigin="anonymous" referrerpolicy="no-referrer" />`
    data_css += `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.41.0/min/vs/editor/editor.main.min.css" integrity="sha512-MFDhxgOYIqLdcYTXw7en/n5BshKoduTitYmX8TkQ+iJOGjrWusRi8+KmfZOrgaDrCjZSotH2d1U1e/Z1KT6nWw==" crossorigin="anonymous" referrerpolicy="no-referrer" />`

    cookie_map := Get_cookie_header(cookies)
    if cookie_map["main_css_darkmode"] == "1" {
        // Main CSS
        data_css_dark += `<link rel="stylesheet" href="/views/main_css/css/sub/dark.css` + data_css_ver + `">`

        // External CSS
        data_css_dark += `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/dark.min.css" integrity="sha512-bfLTSZK4qMP/TWeS1XJAR/VDX0Uhe84nN5YmpKk5x8lMkV0D+LwbuxaJMYTPIV13FzEv4CUOhHoc+xZBDgG9QA==" crossorigin="anonymous" referrerpolicy="no-referrer" />`
    }

    new_data := append([]any{}, data[:2]...)
    new_data = append(new_data, "", data_css)
    new_data = append(new_data, data[2], data_css_dark)
    new_data = append(new_data, data[3:]...)

    log.Default().Println(new_data)

    return new_data
}

func Get_template(db *sql.DB, config Config, name string, data string) string {
    context := pongo2.Context{
        "imp" : []any{
            name,
            Get_wiki_set(db, config.IP, config.Cookies),
            Get_wiki_custom(db, config.IP, config.Session, config.Cookies),
            Get_wiki_css([]any{0, 0}, config.Cookies),
        },
        "data" : data,
        "menu" : 0,
    }

    tpl, err := pongo2.FromFile(Get_skin_route(db, config.IP))
    if err != nil {
        panic(err)
    }

    out, err := tpl.Execute(context)
    if err != nil {
        panic(err)
    }

    return out
}