package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func error_handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                err, ok := r.(error)
                if !ok {
                    err = fmt.Errorf("%v", r)
                }

                stackTrace := debug.Stack()

                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "response" : "error",
                    "error" : err.Error(),
                    "stack" : string(stackTrace),
                })
            }
        }()

        c.Next()
    }
}

func pongo_init() {
    pongo2.RegisterFilter("md5_replace", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
        h := md5.Sum([]byte(in.String()))
        
        return pongo2.AsValue(hex.EncodeToString(h[:])), nil
    })

    pongo2.RegisterFilter("load_lang", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
        db := tool.DB_connect()
        defer tool.DB_close(db)

        return pongo2.AsValue(tool.Get_language(db, in.String(), false)), nil
    })

    pongo2.RegisterFilter("cut_100", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
        s := in.String()
        if len(s) > 100 {
            s = s[:100]
        }

        return pongo2.AsValue(s), nil
    })
}

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
        
    standalone_mode := false

    port := "3001"
    if len(os.Args) > 1 {
        port = os.Args[1]
    }

    var r *gin.Engine
    if len(os.Args) > 2 && os.Args[2] == "dev" {
        r = gin.Default()
    } else {
        gin.SetMode(gin.ReleaseMode)
        r = gin.New()
    }

    if len(os.Args) <= 3 || os.Args[3] != "api" {
        standalone_mode = true
        tool.IN_mod_OUT_mod(standalone_mode)
    }

    r.Use(error_handler())
    pongo_init()
    tool.DB_init()

    r.POST("/compatible_api/:url", func(c *gin.Context) {
        route_data := ""
        
        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            panic(err)
        }
        
        body_string := string(body)

        main_set := map[string]string{}
        main_set["data"] = body_string

        if standalone_mode {
            main_set["ip"] = tool.Get_IP(c)
            main_set["cookies"] = tool.Get_Cookies(c)
            main_set["session"] = ""
        }        

        url_param := c.Param("url")

        config := tool.Config{
            Other_set: main_set["data"],
            IP: main_set["ip"],
            Cookies: main_set["cookies"],
            Session: main_set["session"],
        }
        
        switch url_param {
        case "test":
            route_data = "ok"
        case "api_w_raw":
            route_data = route.Api_w_raw_exter(config)
        case "api_func_sha224":
            route_data = route.Api_func_sha224(config)
        case "api_w_random":
            route_data = route.Api_w_random(config)
        case "api_func_search":
            route_data = route.Api_func_search(config)
        case "api_topic":
            route_data = route.Api_topic(config)
        case "api_func_ip":
            route_data = route.Api_func_ip(config)
        case "api_list_recent_change":
            route_data = route.Api_list_recent_change_exter(config)
        case "api_list_recent_edit_request":
            route_data = route.Api_list_recent_edit_request(config)
        case "api_bbs":
            route_data = route.Api_bbs_exter(config)
        case "api_w_xref":
            route_data = route.Api_w_xref(config)
        case "api_w_watch_list":
            route_data = route.Api_w_watch_list_exter(config)
        case "api_user_watch_list":
            route_data = route.Api_user_watch_list_exter(config)
        case "api_w_render":
            route_data = route.Api_w_render_exter(config)
        case "api_func_llm":
            route_data = route.Api_func_llm(config)
        case "api_func_language":
            route_data = route.Api_func_language(config)
        case "api_func_auth":
            route_data = route.Api_func_auth(config)
        case "api_list_recent_discuss":
            route_data = route.Api_list_recent_discuss(config)
        case "api_bbs_list":
            route_data = route.Api_bbs_list_exter(config)
        case "api_list_old_page":
            route_data = route.Api_list_old_page(config)
        case "api_topic_list":
            route_data = route.Api_topic_list(config)
        case "api_bbs_w_n":
            route_data = route.Api_bbs_w(config)
        case "api_w_set_reset":
            route_data = route.Api_w_set_reset(config)
        case "api_list_recent_block":
            route_data = route.Api_list_recent_block(config)
        case "api_list_title_index":
            route_data = route.Api_list_title_index(config)
        case "api_user_setting_editor_post":
            route_data = route.Api_user_setting_editor_post(config)
        case "api_user_setting_editor_delete":
            route_data = route.Api_user_setting_editor_delete(config)
        case "api_user_setting_editor":
            route_data = route.Api_user_setting_editor(config)
        case "api_setting":
            route_data = route.Api_setting(config)
        case "api_setting_put":
            route_data = route.Api_setting_put(config)
        case "api_func_ip_menu":
            route_data = route.Api_func_ip_menu(config)
        case "api_func_ip_post":
            route_data = route.Api_func_ip_post(config)
        case "api_list_acl":
            route_data = route.Api_list_acl(config)
        case "api_user_rankup":
            route_data = route.Api_user_rankup(config)
        case "api_func_acl":
            route_data = route.Api_func_acl(config)
        case "api_func_ban":
            route_data = route.Api_func_ban(config)
        case "api_func_auth_post":
            route_data = route.Api_func_auth_post(config)
        case "api_give_auth_patch":
            route_data = route.Api_give_auth_patch(config)
        case "api_list_auth":
            route_data = route.Api_list_auth(config)
        case "api_w_page_view":
            route_data = route.Api_w_page_view_exter(config)
        case "api_bbs_w_comment_one":
            route_data = route.Api_bbs_w_comment_one(config, false)
        case "api_bbs_w_comment":
            route_data = route.Api_bbs_w_comment(config)
        case "api_list_history":
            route_data = route.Api_list_history_exter(config)
        case "api_list_markup":
            route_data = route.Api_list_markup(config)
        case "api_bbs_w_set":
            route_data = route.Api_bbs_w_set(config)
        case "api_bbs_w_set_put":
            route_data = route.Api_bbs_w_set_put(config)
        case "api_func_alarm_post":
            route_data = route.Api_func_alarm_post(config)
        case "api_bbs_w":
            route_data = route.Api_bbs_w(config)
        case "api_bbs_w_post":
            route_data = route.Api_bbs_w_post(config)
        case "api_w_comment":
            route_data = route.Api_w_comment(config)
        case "api_bbs_w_tabom":
            route_data = route.Api_bbs_w_tabom(config)
        case "api_bbs_w_tabom_post":
            route_data = route.Api_bbs_w_tabom_post(config)
        case "api_func_email_post":
            route_data = route.Api_func_email_post(config)
        case "api_func_level":
            route_data = route.Api_func_level(config)
        case "api_func_wiki_set":
            route_data = route.Api_func_wiki_set(config)
        case "api_func_skin_name":
            route_data = route.Api_func_skin_name(config)
        case "api_func_wiki_custom":
            route_data = route.Api_func_wiki_custom(config)
        case "api_list_random":
            route_data = route.Api_list_random_exter(config)
        case "view_list_random":
            route_data = route.View_list_random(config).JSON
        case "view_w_watch_list":
            route_data = route.View_w_watch_list(config).JSON
        case "view_user_watch_list":
            route_data = route.View_user_watch_list(config).JSON
        default:
            route_data = "{ \"response\" : \"404\" }"
        }
    
        c.Data(http.StatusOK, "application/json", []byte(route_data))
    })

    r.GET("/list/random", func(c *gin.Context) {
        route_data := route.View_list_random(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })

    r.GET("/random", func(c *gin.Context) {
        route_data := route.View_w_random(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })

    r.GET("/other", func(c *gin.Context) {
        route_data := route.View_main_other(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })

    r.GET("/user", func(c *gin.Context) {
        route_data := route.View_user(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, tool.Get_IP(c)).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })

    r.GET("/user/*user_name", func(c *gin.Context) {
        route_data := route.View_user(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, strings.TrimPrefix(c.Param("user_name"), "/")).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })

    r.GET("/upload", func(c *gin.Context) {
        route_data := route.View_edit_file_upload(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })

    r.GET("/w/*doc_name", func(c *gin.Context) {
        route_data := route.View_w(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, strings.TrimPrefix(c.Param("doc_name"), "/"))
        c.Data(route_data.ST, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.GET("/recent_change", func(c *gin.Context) {
        route_data := route.View_list_recent_change(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, "", "50", "1")
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.GET("/recent_changes", func(c *gin.Context) {
        route_data := route.View_list_recent_change(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, "", "50", "1")
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.GET("/recent_change/:num/:set_type", func(c *gin.Context) {
        route_data := route.View_list_recent_change(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, c.Param("set_type"), "50", c.Param("num"))
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.GET("/bbs/main", func(c *gin.Context) {
        route_data := route.View_bbs_main(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, "1")
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.GET("/history/*doc_name", func(c *gin.Context) {
        route_data := route.View_history(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.GET("/edit/*doc_name", func(c *gin.Context) {
        route_data := route.View_edit(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, strings.TrimPrefix(c.Param("doc_name"), "/"))
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.POST("/edit/*doc_name", func(c *gin.Context) {
        doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
        data := c.PostForm("content")
        send := c.PostForm("send")
        agree := c.PostForm("copyright_agreement")

        route_data := route.View_edit_post(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, doc_name, data, send, agree)
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data.HTML))
    })

    r.POST("/upload", func(c *gin.Context) {
        form, err := c.MultipartForm()
        if err != nil || form == nil {
            c.String(http.StatusBadRequest, "invalid multipart form")
            return
        }

        files := form.File["f_data[]"]
        if len(files) == 0 {
            c.String(http.StatusBadRequest, "no file")
            return
        }

        posted_name := strings.TrimSpace(c.PostForm("f_name"))
        other_set_arr := []map[string]string{}

        count := 1
        for _, fh := range files {
            f, err := fh.Open()
            if err != nil {
                continue
            }
            
            b, err := io.ReadAll(f)
            
            _ = f.Close()
            if err != nil {
                continue
            }

            name := posted_name

            name = strings.TrimSpace(name)
            ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
            ext = strings.TrimSpace(ext)

            b64 := base64.StdEncoding.EncodeToString(b)

            other_set := map[string]string{
                "file_name": name,
                "file_ext": ext,
                "file_data": b64,
            }

            other_set_arr = append(other_set_arr, other_set)
            count += 1
        }

        other_set_arr_str, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(other_set_arr)

        route_data := route.View_edit_file_upload_post(tool.Config{
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
            Other_set: other_set_arr_str,
        }).HTML
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
    })
    
    r.GET("/view/*name", route.View_view_file)
    r.GET("/views/*name", route.View_view_file)
    r.GET("/image/*name", route.View_view_image_file)

    r.NoRoute(func(c *gin.Context) {
        route_data := route.View_main_404_page(tool.Config{
            Other_set: "",
            IP: tool.Get_IP(c),
            Cookies: tool.Get_Cookies(c),
            Session: "",
        }, c.Request.URL.Path).HTML
        c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(route_data))
    })

    if standalone_mode {
        log.Default().Println("Run in http://localhost:" + port)
        r.Run("0.0.0.0:" + port)
    } else {
        r.Run("127.0.0.1:" + port)
    }
}
