package servers

import (
	"github.com/IzePhanthakarn/kawaii-shop/modules/files/filesUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/modules/products/productsHandlers"
	"github.com/IzePhanthakarn/kawaii-shop/modules/products/productsRepositories"
	"github.com/IzePhanthakarn/kawaii-shop/modules/products/productsUsecases"
)

type IProductsModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type productsModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductsModule {
	fileUsecase := filesUsecases.FileUsecase(m.server.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.server.db, m.server.cfg, fileUsecase)
	productsUsecase := productsUsecases.ProductsUsecase(productsRepository)
	productsHandler := productsHandlers.ProductsHandler(m.server.cfg, fileUsecase, productsUsecase)

	return &productsModule{
		moduleFactory: m,
		repository:    productsRepository,
		usecase:       productsUsecase,
		handler:       productsHandler,
	}
}

func (p *productsModule) Init() {
	router := p.router.Group("/products")

	router.Get("/", p.handler.FindProduct)
	router.Get("/:product_id", p.handler.FindOneProduct)

	router.Post("/", p.handler.AddProduct, p.middlewares.JwtAuth(), p.middlewares.Authorize(2))
	router.Patch("/:product_id", p.handler.UpdateProduct, p.middlewares.JwtAuth(), p.middlewares.Authorize(2))

	router.Delete("/:product_id", p.handler.DeleteProduct, p.middlewares.JwtAuth(), p.middlewares.Authorize(2))
}

func (f *productsModule) Repository() productsRepositories.IProductsRepository { return f.repository }

func (f *productsModule) Usecase() productsUsecases.IProductsUsecase { return f.usecase }

func (f *productsModule) Handler() productsHandlers.IProductsHandler { return f.handler }