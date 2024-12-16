package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func sendRequestAndGetResponse(apiKey string, redisServer *RedisServer) (string, error) {
    // 构建请求数据
    data := map[string]interface{}{
        "model": "glm-4-flash",
        "messages": redisServer.SessionList, // 使用 RedisServer 结构体的会话列表
    }
    fmt.Println(data)

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
    fmt.Println(response)
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
