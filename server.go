// @Title  server.go
// @Description High Interaction Honeypot Solution for Redis protocol
// @Author  Cy 2021.04.08
package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/sirupsen/logrus"
	"github.com/walu/resp"
	"gopkg.in/ini.v1"
)

type SessionMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RedisServer struct {
    server  *gev.Server
    hashmap *hashmap.Map
    Config  *ini.File
    log     *logrus.Logger
    SessionList []SessionMessage
}

func NewRedisServer(address string, proto string, loopsnum int) (server *RedisServer, err error) {
	Serv := new(RedisServer)
	Serv.hashmap = hashmap.New()
	config, err := LoadConfig("redis.conf")
	Serv.log = logrus.New()
	Serv.log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if err != nil {
		panic(err)
	}
	Serv.Config = config
	Serv.server, err = gev.NewServer(Serv,
		gev.Address(address),
		gev.Network(proto),
		gev.NumLoops(loopsnum))
	if err != nil {
		return nil, err
	}

	// init llm info
	updateSessionList(Serv,"请你扮演一个 redis 服务器，我将会以 redis 客户端的身份，通过命令行的形式与你沟通，你需要模仿真实 redis 服务器能够给出的响应向我回复。请注意，我只需要你给出命令响应，不需要任何其他的解释或分析。","system")
	return Serv, nil
}

func (s *RedisServer) Start() {
	s.server.Start()
}

func (s *RedisServer) Stop() {
	s.server.Stop()
}

func (s *RedisServer) OnConnect(c *connection.Connection) {
	s.log.WithFields(logrus.Fields{
		"action": "NewConnect",
		"addr":   c.PeerAddr(),
	}).Println()
}

func (s *RedisServer) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	command := bytes.NewReader(data)
	if command.Len() == 2 {
		return
	}
	cmd, err := resp.ReadCommand(command)
	if err != nil {
		return
	}

	com := strings.ToLower(cmd.Name())
	action := strings.Join(cmd.Args, " ")

	s.log.WithFields(logrus.Fields{
		"action": action,
		"addr":   c.PeerAddr(),
	}).Println()

	switch com {
	case "ping":
		out = []byte("+PONG\r\n")
	case "flushall":
		out = []byte("+OK\r\n")
	case "flushdb":
		out = []byte("+OK\r\n")
	case "save":
		out = []byte("+OK\r\n")
	case "select":
		out = []byte("+OK\r\n")
	default:
		apiKey := ""
		updateSessionList(s, action, "user")
		response, err := sendRequestAndGetResponse(apiKey,s)
		if err!= nil {
			fmt.Println("请求出错:", err)
			return
		}
		extract_content, err := extractMessageContent(response)
		if err != nil {
			fmt.Println("json解析出错", err)
			return
		}

		s.log.WithFields(logrus.Fields{
			"llm response": extract_content,
			"addr":   c.PeerAddr(),
		}).Println()

		out = []byte("+"+extract_content+"\r\n")
		updateSessionList(s, extract_content, "assistant")
	}
	return
}

func updateSessionList(rs *RedisServer, message string,role string) {
    sessionMessage := SessionMessage{
            Role:    role,
            Content: message,
        }
        rs.SessionList = append(rs.SessionList, sessionMessage)
}

func (s *RedisServer) OnClose(c *connection.Connection) {
	s.log.WithFields(logrus.Fields{
		"action": "Closed",
		"addr":   c.PeerAddr(),
	}).Println()
}
