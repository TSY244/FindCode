package app

import (
	"ScanIDOR/internal/pkg/server/config"
	"ScanIDOR/internal/pkg/server/models"
	"ScanIDOR/internal/pkg/server/repositories"
	"ScanIDOR/internal/pkg/server/router"
	"ScanIDOR/internal/pkg/server/services"
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Config     *config.Config
	router     *gin.Engine
	db         *gorm.DB
	RootCtx    context.Context
	RootCancel context.CancelFunc
}

func NewApp(c *config.Config) (*App, error) {
	// init logger

	app := App{
		Config: c,
	}

	app.initDatabase(c)

	// 创建service
	userService := services.NewUserService(repositories.NewUserRepo(app.db), repositories.NewTokenRepo(app.db))
	findCodeService := services.NewFindCodeService()
	app.router = router.SetupRouter(findCodeService, userService)

	return &app, nil
}

func (s *App) initDatabase(c *config.Config) {
	db, err := gorm.Open(c.DbConfig.GetDrive(), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Token{},
	); err != nil {
		panic(err)
	}
	s.db = db

}

func (s *App) Run() {
	s.router.LoadHTMLGlob("web/templates/*")
	s.router.Static("/asset", "web/asset")
	//s.router.Static("/asset/js", "web/asset/js")

	s.router.Run(":8080")
}
