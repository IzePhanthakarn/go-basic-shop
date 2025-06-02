package monitorHandlers

import (
	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/modules/monitor"
	fiber "github.com/gofiber/fiber/v3"
)

type IMonitorHandler interface {
	HealthCheck(c fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

// @Summary HealthCheck
// @Description HealthCheck
// @Tags Monitor
// @Accept  json
// @Produce  json
// @Success 200
// @Router / [get]
func (m *monitorHandler) HealthCheck(c fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    m.cfg.App().Name(),
		Version: m.cfg.App().Version(),
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
