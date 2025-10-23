package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func View_edit_post(config tool.Config, doc_name string, data string, send string, agree string) tool.View_result {   
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    return_data := Api_edit_post(config, doc_name, data, send, agree)

    result_html := ""
    if return_data["response"].(string) == "ok" {
        result_html = tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
    } else {
        result_html = return_data["data"].(string)
    }

    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : result_html,
        JSON : string(json_data),
    }

    return result_data
}