package http

import (
	"time"

	"github.com/SoundBoardBot/server-counter/config"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	Logger *zap.Logger
	Config config.Config
	router *gin.Engine
}

func NewServer(logger *zap.Logger, config config.Config) *Server {
	return &Server{
		Logger: logger,
		Config: config,
		router: gin.New(),
	}
}

func (s *Server) RegisterRoutes() {
	s.router.Use(ginzap.Ginzap(s.Logger, time.RFC3339, true))
	s.router.Use(ginzap.RecoveryWithZap(s.Logger, true))

	s.router.GET("/metrics", s.metricsGetHandler)
}

func (s *Server) Start() {
	if err := s.router.Run(s.Config.Http.Address); err != nil {
		panic(err)
	}
}
