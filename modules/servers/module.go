package servers

import (
	"github.com/Flussen/swagger-fiber-v3"
	"github.com/IzePhanthakarn/go-basic-shop/modules/appinfo/appinfoHandlers"
	"github.com/IzePhanthakarn/go-basic-shop/modules/appinfo/appinfoRepositories"
	"github.com/IzePhanthakarn/go-basic-shop/modules/appinfo/appinfoUsecases"
	"github.com/IzePhanthakarn/go-basic-shop/modules/files/filesUsecases"
	"github.com/IzePhanthakarn/go-basic-shop/modules/middlewares/middlewaresHandlers"
	"github.com/IzePhanthakarn/go-basic-shop/modules/middlewares/middlewaresRepositories"
	"github.com/IzePhanthakarn/go-basic-shop/modules/middlewares/middlewaresUsecases"
	"github.com/IzePhanthakarn/go-basic-shop/modules/monitor/monitorHandlers"
	"github.com/IzePhanthakarn/go-basic-shop/modules/orders/ordersHandlers"
	"github.com/IzePhanthakarn/go-basic-shop/modules/orders/ordersRepositories"
	"github.com/IzePhanthakarn/go-basic-shop/modules/orders/ordersUsecases"
	"github.com/IzePhanthakarn/go-basic-shop/modules/products/productsRepositories"
	"github.com/IzePhanthakarn/go-basic-shop/modules/users/usersHandlers"
	"github.com/IzePhanthakarn/go-basic-shop/modules/users/usersRepositories"
	"github.com/IzePhanthakarn/go-basic-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
	AppinfoModule()
	FileModule() IFileModule
	ProductsModule() IProductsModule
	OrderModule()
	SwaggerModule()
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
		middlewares: middlewares,
	}
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

func (m *moduleFactory) SwaggerModule() {
	m.router.Get("/swagger/*", swagger.HandlerDefault)
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
	router.Post("/signup-admin", handler.SignUpAdmin, m.middlewares.JwtAuth(), m.middlewares.Authorize(2))

	router.Get("/:user_id", handler.GetUserProfile, m.middlewares.JwtAuth(), m.middlewares.ParamsCheck())
	router.Get("/admin/secret", handler.GenerateAdminToken, m.middlewares.JwtAuth(), m.middlewares.Authorize(2))
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Post("/categories", handler.AddCategory, m.middlewares.JwtAuth(), m.middlewares.Authorize(2))
	router.Delete("/categories/:category_id", handler.RemoveCategory, m.middlewares.JwtAuth(), m.middlewares.Authorize(2))

	router.Get("/categories", handler.FindCategory)
	router.Get("/apikey", handler.GenerateApiKey, m.middlewares.JwtAuth(), m.middlewares.Authorize(2))
}

func (m *moduleFactory) OrderModule() {
	filesUsecase := filesUsecases.FileUsecase(m.server.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.server.db, m.server.cfg, filesUsecase)

	ordersRepository := ordersRepositories.OrdersRepository(m.server.db)
	ordersUsecase := ordersUsecases.OrderUsecase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrdersHandlers(m.server.cfg, ordersUsecase)

	router := m.router.Group("/orders")

	router.Post("/", ordersHandler.InsertOrder, m.middlewares.JwtAuth())

	router.Get("/", ordersHandler.FindOrder, m.middlewares.JwtAuth())
	router.Get("/:user_id/:order_id", ordersHandler.FindOneOrder, m.middlewares.JwtAuth(), m.middlewares.ParamsCheck())

	router.Patch("/:user_id/:order_id", ordersHandler.UpdateOrder, m.middlewares.JwtAuth(), m.middlewares.ParamsCheck())
}
