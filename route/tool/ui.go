package tool

import (
	"database/sql"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/flosch/pongo2/v6"
)

func Get_skin_route(skin_name string, route string) string {
    return filepath.Join("..", "views", skin_name, route)
}

func Get_template_set(skin_name string) map[string]string {
    set_file_path := Get_skin_route(skin_name, "set.json")
	if _, err := os.Stat(set_file_path); err == nil {
		data, err := os.ReadFile(set_file_path)
		if err != nil {
			panic(err)
		}

		set_json := map[string]string{}
    	json.Unmarshal([]byte(data), &set_json)

		return set_json
	}

    return map[string]string{}
}

func Get_template(db *sql.DB, config Config, name string, data string, sub string, menu [][]any) string {
    skin_name := Get_use_skin_name(db, config.IP)
    
    template_set := Get_template_set(skin_name)
    for k, v := range template_set {
        data = strings.ReplaceAll(data, k, v)
    }

	context := pongo2.Context{
		"imp": []any{
			name,
			Get_wiki_set(db, config.IP, config.Cookies),
			Get_wiki_custom(db, config.IP, config.Session, config.Cookies),
			Get_wiki_css([]any{func(sub string) any {
				if sub == "" {
					return 0
				} else {
					return sub
				}
			}(sub), 0}, config.Cookies),
		},
		"data": `<div class="opennamu_main">` + data + `</div>`,
		"menu": func(menu [][]any) any {
			if len(menu) == 0 {
				return 0
			} else {
				return menu
			}
		}(menu),
	}

	tpl, err := pongo2.FromFile(Get_skin_route(skin_name, "index.html"))
	if err != nil {
		panic(err)
	}

	out, err := tpl.Execute(context)
	if err != nil {
		panic(err)
	}

	return out
}

func Get_redirect(target string) string {
    attrURL := html.EscapeString(target)
    jsURL := strconv.Quote(target)

    return fmt.Sprintf(`<!doctype html>
<html lang="ko">
<head>
<meta charset="utf-8">
<title>Redirecting…</title>
<script>
location.replace(%s);
</script>
<noscript>
<meta http-equiv="refresh" content="0; url=%s">
</noscript>
</head>
<body>
<p>Redirecting… <a href="%s">continue</a></p>
</body>
</html>`, jsURL, attrURL, attrURL)
}

type View_result struct {
    HTML string
    JSON string
    ST int
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

    end := 2
    if end > len(data) {
        end = len(data)
    }

    new_data := append([]any{}, data[:end]...)
    new_data = append(new_data, "", data_css)

    if len(data) >= 3 {
        new_data = append(new_data, data[2])
    }

    new_data = append(new_data, data_css_dark)

    if len(data) >= 3 {
        new_data = append(new_data, data[3:]...)
    }

    return new_data
}

func Get_list_ui(left string, right string, bottom string, class_name string) string {
    data_html := ""

    data_html += `<span class="` + class_name + `">`
    data_html += `<div class="opennamu_recent_change">`
    data_html += left

    data_html += `<div style="float: right;">`
    data_html += right
    data_html += `</div>`

    data_html += `<div style="clear: both;"></div>`

    if bottom != "" {
        data_html += "<hr>"
        data_html += bottom
    }

    data_html += "</div>"
    data_html += `<hr class="main_hr">`
    data_html += "</span>"

    return data_html
}

func Get_error_page(db *sql.DB, config Config, error_name string) string {
    data := ""
    if error_name == "auth" {
        data = Get_language(db, "authority_error", true)
    }

    return Get_template(
        db,
        config,
        Get_language(db, "error", true),
        `<h2>` + Get_language(db, "error", true) + `</h2>` +
        `<ul>` +
            `<li>` + data + `</li>` +
        `</ul>`,
        "",
        [][]any{},
    )
}

func Get_page_control(db *sql.DB, page int, count int, max_count int, url string) string {
    data_html := "<hr class=\"main_hr\">"

    if page > 1 {
        prev_page := page - 1
        before_url := strings.ReplaceAll(url, "{}", strconv.Itoa(prev_page))

        data_html += `<a href="` + before_url + `">(` + Get_language(db, "previous", true) + `)</a> `
    } else if count == max_count {
        prev_page := page + 1
        after_url := strings.ReplaceAll(url, "{}", strconv.Itoa(prev_page))

        data_html += `<a href="` + after_url + `">(` + Get_language(db, "next", true) + `)</a> `
    }

    return data_html
}