package servers

import (
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresHandlers"
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresRepositories"
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/modules/monitor/monitorHandlers"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users/usersHandlers"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users/usersRepositories"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
}

type moduleFactory struct {
	router      fiber.Router
	server      *server
	middlewares middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(router fiber.Router, server *server, middlewares middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router:      router,
		server:      server,
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

func (m *moduleFactory) UserModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)
	
	router.Get("/secret", handler.GenerateAdminToken)
}
