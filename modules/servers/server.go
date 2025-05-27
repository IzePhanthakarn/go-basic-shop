package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type Iserver interface {
	Start()
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) Iserver {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	// Middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	// Modules
	v1 := s.app.Group("/v1")
	modules := InitModule(v1, s, middlewares)

	modules.MonitorModule()
	modules.UserModule()
	modules.AppinfoModule()

	s.app.Use(middlewares.RouterCheck())

	// Greaceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	log.Println("server is running on", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
