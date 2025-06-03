package servers

import (
	"github.com/IzePhanthakarn/go-basic-shop/modules/files/filesHandlers"
	"github.com/IzePhanthakarn/go-basic-shop/modules/files/filesUsecases"
)

type IFileModule interface {
	Init()
	Usecase() filesUsecases.IFilesUsecase
	Handler() filesHandlers.IFilesHandler
}

type filesModule struct {
	*moduleFactory
	usecase filesUsecases.IFilesUsecase
	handler filesHandlers.IFilesHandler
}

func (m *moduleFactory) FileModule() IFileModule {
	usecase := filesUsecases.FileUsecase(m.server.cfg)
	handler := filesHandlers.FilesHandler(m.server.cfg, usecase)

	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *filesModule) Init() {
	router := f.router.Group("/files")
	router.Post("/upload", f.handler.UploadFile, f.middlewares.JwtAuth(), f.middlewares.Authorize(2))
	router.Patch("/delete", f.handler.DeleteFile, f.middlewares.JwtAuth(), f.middlewares.Authorize(2))
}

func (f *filesModule) Usecase() filesUsecases.IFilesUsecase { return f.usecase }

func (f *filesModule) Handler() filesHandlers.IFilesHandler { return f.handler }
