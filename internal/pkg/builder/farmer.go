package builder

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"farmer/internal/pkg/api"
	cfg "farmer/internal/pkg/config"
	spotmanager "farmer/internal/pkg/spot_manager"
	"farmer/internal/pkg/telebot"
)

type ISpotFarmerSystem interface {
	Run() error
}

type spotFarmerSystem struct {
	server *gin.Engine
	m      spotmanager.ISpotManager
	t      telebot.ITeleBot
}

var sys *spotFarmerSystem

func NewSpotFarmerSystem(m spotmanager.ISpotManager, t telebot.ITeleBot) (ISpotFarmerSystem, error) {
	server, err := newServer()
	if err != nil {
		return nil, errors.New("can not build gin server")
	}
	api.AddRouterV1(server)

	sys = &spotFarmerSystem{
		server: server,
		m:      m,
		t:      t,
	}
	return sys, nil
}

func (sys *spotFarmerSystem) Run() error {
	go sys.m.Run()
	go sys.t.Run()

	fmt.Println(cfg.Instance().Http.BindAddress)
	return sys.server.Run(cfg.Instance().Http.BindAddress)
}

func newServer() (*gin.Engine, error) {
	gin.SetMode(cfg.Instance().Http.Mode)
	server := gin.Default()
	setCORS(server)
	return server, nil
}

func setCORS(engine *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowMethods(http.MethodOptions)
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))
}
