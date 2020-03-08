package rest

import (
	"api-crawler/core/collegescorecard/collegeapplication"
	"api-crawler/infrastructure/db/sqliteinfra"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sevenNt/echo-pprof"
	"log"
)

type Config struct {
	ApiKey     string `env:"API_KEY" envDefault:"tzj89Y0IAMihf24Qro9CgiUiHAGFAntQPWf08Rvb"`
	CollegeUrl string `env:"COLLEGE_URL" envDefault:"https://api.data.gov/ed/collegescorecard/v1/schools.json?"`
	ServerPort string `env:"SERVER_PORT" envDefault:":8080"`
}

type service struct {
	collegeService collegeapplication.CollegeService
	echoServer     *echo.Echo
	config         *Config
	db             *gorm.DB
}

func NewService(db *gorm.DB) *service {
	//config
	config := new(Config)
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
	//setUp application services
	collegeRepository := sqliteinfra.NewCollegeRepository(db)
	collegeService := collegeapplication.NewCollegeService(config.ApiKey, config.CollegeUrl, collegeRepository)

	//setupEcho
	echoServer := echo.New()
	echopprof.Wrap(echoServer)

	service := &service{
		collegeService: collegeService,
		echoServer:     echoServer,
		config:         config,
		db:             db,
	}
	return service
}

func (s *service) setUpMiddleware() {
	// Middleware
	s.echoServer.Use(middleware.Logger())
	s.echoServer.Use(middleware.Recover())
}

func (s *service) setUpRoute() {
	// Routes
	s.echoServer.GET("/", s.hello)
	s.echoServer.GET("/schools", s.getAllSchools)
	s.echoServer.GET("/num", s.num)
}

func (s *service) StartServer() {
	s.setUpMiddleware()
	s.setUpRoute()
	go func() {
		err := s.collegeService.CrawlApi()
		if err != nil {
			panic(err)
		}
	}()
	s.echoServer.Logger.Fatal(s.echoServer.Start(s.config.ServerPort))

}
