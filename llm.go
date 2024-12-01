package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func sendRequestAndGetResponse(apiKey, content string) (string, error) {
    // 构建请求数据
    data := map[string]interface{}{
        "model": "glm-4-flash",
        "messages": []map[string]string{
            {
                "role": "system",
                "content": "请你扮演一个 redis 服务器，我将会以 redis 客户端的身份，通过命令行的形式与你沟通，你需要模仿真实 redis 服务器能够给出的响应向我回复。请注意，我只需要你给出命令响应，不需要任何其他的解释或分析。",
            },
            {
                "role": "user",
                "content": content,
            },
        },
    }

    jsonData, err := json.Marshal(data)
    if err!= nil {
        return "", err
    }

    // 发送请求
    url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"
    req, err := http.NewRequest("POST", url, nil)
    if err!= nil {
        return "", err
    }

    req.Header.Add("Authorization", "Bearer "+apiKey)
    req.Header.Add("Content-Type", "application/json")

    req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))

    client := &http.Client{}
    res, err := client.Do(req)
    if err!= nil {
        return "", err
    }
    defer res.Body.Close()

    // 读取响应内容
    body, err := ioutil.ReadAll(res.Body)
    if err!= nil {
        return "", err
    }

    return string(body), nil
}

func extractMessageContent(response string) (string, error) {
    var data map[string]interface{}
    json.Unmarshal([]byte(response), &data)

    choices := data["choices"].([]interface{})
    if len(choices) > 0 {
        choice := choices[0].(map[string]interface{})
        message := choice["message"].(map[string]interface{})
        return message["content"].(string),nil
    }

    return "",nil
}

// func main() {
//     apiKey := ""
//     content := "ping"

//     response, err := sendRequestAndGetResponse(apiKey, content)
//     if err!= nil {
//         fmt.Println("请求出错:", err)
//         return
//     }
// 	extract_content, err := extractMessageContent(response) 
// 	if err != nil {
// 		fmt.Println("json解析出错", err)
// 		return 
// 	}

//     fmt.Println(extract_content)
// }
