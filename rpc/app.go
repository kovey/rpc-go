package rpc

import (
	"github.com/kovey/config-go/config"
	"github.com/kovey/logger-go/logger"
	"github.com/kovey/logger-go/monitor"
	"github.com/kovey/rpc-go/router"
	"github.com/kovey/server-go/server"
)

var (
	app  *App
	conf *config.Config
)

type App struct {
	server   *server.Server
	config   *server.Config
	event    RpcEvent
	confPath string
}

func NewApp(conf *server.Config, confPath string) *App {
	return &App{event: NewRpcEvent(), config: conf, confPath: confPath}
}

func Router(path string, method string, handler interface{}) {
	router.Route.Add(path, router.NewRouter(path, method, handler))
	return
}

func RouterByHandler(path string, handler interface{}) {
	router.Route.AddAll(path, handler)
}

func Get(path string, method string) (*router.Router, error) {
	return router.Route.Get(path, method)
}

func IsExists(path string, method string) bool {
	return router.Route.IsExists(path, method)
}

func Run(confPath string, before func(app *App, config *config.Config), after func(app *App, config *config.Config)) {
	var err error
	conf, err = config.LoadConfig(confPath)
	if err != nil {
		panic(err)
	}

	servConf := server.NewConfig(conf.Server.Package.OpenCheck, conf.Server.Package.BodyLenOffset, conf.Server.Package.HeaderLen, server.INT_32)
	app = NewApp(servConf, confPath)
	logger.SetLevelByName(conf.Logger.Level)
	logger.GetInstance().SetDir(conf.Logger.LogDir).SetLevelByName(conf.Logger.Level)
	monitor.Init(conf.Logger.LogDir + "/monitor")

	app.server = server.NewServer(conf.Server.Host, conf.Server.Port, server.SOCKET_TCP)
	app.server.Set(app.config)
	app.server.SetEvent(app.event)

	if before != nil {
		before(app, conf)
	}

	app.server.Start()

	if after != nil {
		after(app, conf)
	}
}
