package route

import (
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func View_edit_post(config tool.Config, c *gin.Context) tool.View_result {   
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    
    doc_name := c.Param("doc_name")
	data := c.PostForm("content")
    send := c.PostForm("send")
    agree := c.PostForm("copyright_agreement")

    return_data := Api_edit_post(config, doc_name, data, send, agree)

    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : tool.Get_redirect("/w/" + tool.Url_parser(doc_name)),
        JSON : string(json_data),
    }

    return result_data
}