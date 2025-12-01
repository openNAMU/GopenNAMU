package route

import (
	"opennamu/route/tool"
)

func View_edit_file_upload_post(config tool.Config) tool.View_result {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := []map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    resp := []any{}
    for _, v := range other_set {
        var config_sub tool.Config

        b, err := json.Marshal(config)
        if err != nil {
            continue
        }

        err = json.Unmarshal(b, &config_sub)
        if err != nil {
            continue
        }

        v_str, _ := json.MarshalToString(v)
        config_sub.Other_set = v_str

        data := Api_file_upload_post(config_sub)

        data_sub := map[string]any{}
        json.Unmarshal([]byte(data), &data_sub)

        resp = append(resp, data_sub["data"])
    }

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = resp
    
    json_data, _ := json.Marshal(return_data)

    result_data := tool.View_result{
        HTML : "",
        JSON : string(json_data),
    }

    return result_data
}