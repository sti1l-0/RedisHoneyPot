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

type RedisServer struct {
	server  *gev.Server
	hashmap *hashmap.Map
	Config  *ini.File
	log     *logrus.Logger
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
		panic(err)
	}
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

	s.log.WithFields(logrus.Fields{
		"action": strings.Join(cmd.Args, " "),
		"addr":   c.PeerAddr(),
	}).Println()

	apiKey := ""
	response,err := sendRequestAndGetResponse(apiKey, com)
    if err!= nil {
        fmt.Println("请求出错:", err)
        return
    }
	extract_content, err := extractMessageContent(response) 
	if err != nil {
		fmt.Println("json解析出错", err)
		return 
	}
	out = []byte(extract_content)
	return	out
}

func (s *RedisServer) OnClose(c *connection.Connection) {
	s.log.WithFields(logrus.Fields{
		"action": "Closed",
		"addr":   c.PeerAddr(),
	}).Println()
}
