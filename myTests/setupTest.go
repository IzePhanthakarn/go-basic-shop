package myTests

import (
	"encoding/json"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/servers"
	"github.com/IzePhanthakarn/kawaii-shop/pkg/databases"

	_ "github.com/IzePhanthakarn/kawaii-shop/docs"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnect(cfg.Db())

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}