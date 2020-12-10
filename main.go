package main

import (
	log "github.com/sirupsen/logrus"

	"io/ioutil"

	"encoding/json"

	"github.com/0B1t322/tasker/tasker_server"
	"google.golang.org/grpc"
)

var conf taskerserver.Config

func init() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}
}

func main() {
	log.Info("Start tasker on " + conf.Port)
	if err := taskerserver.NewTaskerServer(conf, []grpc.ServerOption{}); err != nil {
		panic(err)
	}
}