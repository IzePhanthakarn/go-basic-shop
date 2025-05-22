package servers

import (
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresHandlers"
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresRepositories"
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	middlewares middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(router fiber.Router, server *server, middlewares middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router:     router,
		server:     server,
		middlewares: middlewares}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	handler := middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
	return handler
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}
