package server

import (
	"ScanIDOR/internal/pkg/server/router"
	"github.com/gin-gonic/gin"
)

type Server struct {
	ginServer *gin.Engine
}

func NewServer() *Server {
	return &Server{gin.Default()}
}

func (s *Server) Run() {
	s.ginServer.LoadHTMLGlob("templates/*")
	router.SetupRouter(s.ginServer)
	s.ginServer.Run(":8080")
}
