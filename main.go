package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"opennamu/route"
	"opennamu/route/tool"

	"net/http"

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
                log.Default().Printf("Recovered from panic: %v\nStack Trace:\n%s", err, stackTrace)

                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "response" : "error",
                })
            }
        }()

        c.Next()
    }
}

func main() {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    log.SetFlags(log.LstdFlags | log.Lshortfile)
        
    var r *gin.Engine
    if len(os.Args) > 1 && os.Args[1] == "dev" {
        r = gin.Default()
    } else {
        gin.SetMode(gin.ReleaseMode)
        r = gin.New()
    }

    r.Use(error_handler())

    r.POST("/", func(c *gin.Context) {
        route_data := ""
        
        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            panic(err)
        }
        
        body_string := string(body)

        main_set := map[string]string{}
        json.Unmarshal([]byte(body_string), &main_set)

        if main_set["url"] == "test" {
            tool.DB_init(main_set["data"])
        }

        db := tool.DB_connect()
        defer tool.DB_close(db)

        if len(os.Args) > 1 && os.Args[1] == "dev" {
            log.Default().Println(main_set["url"])
        }

        config := tool.Config{
            Other_set: main_set["data"],
            IP: main_set["ip"],
            Cookies: main_set["cookies"],
            Session: main_set["session"],
        }
        
        switch main_set["url"] {
        case "test":
            route_data = "ok"
        case "main_func_easter_egg":
            route_data = route.View_main_func_easter_egg()
        case "api_w_raw":
            route_data = route.Api_w_raw(db, config)
        case "api_func_sha224":
            route_data = route.Api_func_sha224(db, config)
        case "api_w_random":
            route_data = route.Api_w_random(db, config)
        case "api_func_search":
            route_data = route.Api_func_search(db, config)
        case "api_topic":
            route_data = route.Api_topic(db, config)
        case "api_func_ip":
            route_data = route.Api_func_ip(db, config)
        case "api_list_recent_change":
            route_data = route.Api_list_recent_change(db, config)
        case "api_list_recent_edit_request":
            route_data = route.Api_list_recent_edit_request(db, config)
        case "api_bbs":
            route_data = route.Api_bbs(db, config)
        case "api_w_xref":
            route_data = route.Api_w_xref(db, config)
        case "api_w_watch_list":
            route_data = route.Api_w_watch_list(db, config)
        case "api_user_watch_list":
            route_data = route.Api_user_watch_list(db, config)
        case "api_w_render":
            route_data = route.Api_w_render(db, config)
        case "api_func_llm":
            route_data = route.Api_func_llm(db, config)
        case "api_func_language":
            route_data = route.Api_func_language(db, config)
        case "api_func_auth":
            route_data = route.Api_func_auth(db, config)
        case "api_list_recent_discuss":
            route_data = route.Api_list_recent_discuss(db, config)
        case "api_bbs_list":
            route_data = route.Api_bbs_list(db, config)
        case "api_list_old_page":
            route_data = route.Api_list_old_page(db, config)
        case "api_topic_list":
            route_data = route.Api_topic_list(db, config)
        case "api_bbs_w_n":
            route_data = route.Api_bbs_w(db, config)
        case "api_w_set_reset":
            route_data = route.Api_w_set_reset(db, config)
        case "api_list_recent_block":
            route_data = route.Api_list_recent_block(db, config)
        case "api_list_title_index":
            route_data = route.Api_list_title_index(db, config)
        case "api_user_setting_editor_post":
            route_data = route.Api_user_setting_editor_post(db, config)
        case "api_user_setting_editor_delete":
            route_data = route.Api_user_setting_editor_delete(db, config)
        case "api_user_setting_editor":
            route_data = route.Api_user_setting_editor(db, config)
        case "api_setting":
            route_data = route.Api_setting(db, config)
        case "api_setting_put":
            route_data = route.Api_setting_put(db, config)
        case "api_func_ip_menu":
            route_data = route.Api_func_ip_menu(db, config)
        case "api_func_ip_post":
            route_data = route.Api_func_ip_post(db, config)
        case "api_list_acl":
            route_data = route.Api_list_acl(db, config)
        case "api_user_rankup":
            route_data = route.Api_user_rankup(db, config)
        case "api_func_acl":
            route_data = route.Api_func_acl(db, config)
        case "api_func_ban":
            route_data = route.Api_func_ban(db, config)
        case "api_func_auth_post":
            route_data = route.Api_func_auth_post(db, config)
        case "api_give_auth_patch":
            route_data = route.Api_give_auth_patch(db, config)
        case "api_list_auth":
            route_data = route.Api_list_auth(db, config)
        case "api_w_page_view":
            route_data = route.Api_w_page_view(db, config)
        case "api_bbs_w_comment_one":
            route_data = route.Api_bbs_w_comment_one(db, config, false)
        case "api_bbs_w_comment":
            route_data = route.Api_bbs_w_comment(db, config)
        case "api_list_history":
            route_data = route.Api_list_history(db, config)
        case "api_list_markup":
            route_data = route.Api_list_markup(db, config)
        case "api_bbs_w_set":
            route_data = route.Api_bbs_w_set(db, config)
        case "api_bbs_w_set_put":
            route_data = route.Api_bbs_w_set_put(db, config)
        case "api_func_alarm_post":
            route_data = route.Api_func_alarm_post(db, config)
        case "api_bbs_w":
            route_data = route.Api_bbs_w(db, config)
        case "api_bbs_w_post":
            route_data = route.Api_bbs_w_post(db, config)
        case "api_w_comment":
            route_data = route.Api_w_comment(db, config)
        case "api_bbs_w_tabom":
            route_data = route.Api_bbs_w_tabom(db, config)
        case "api_bbs_w_tabom_post":
            route_data = route.Api_bbs_w_tabom_post(db, config)
        case "api_func_email_post":
            route_data = route.Api_func_email_post(db, config)
        default:
            route_data = "{ \"response\" : \"404\" }"
        }
    
        c.Data(http.StatusOK, "application/json", []byte(route_data))
    })
    
    r.Run(":" + tool.Get_port())
}
