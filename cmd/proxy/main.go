package main

import (
	"sync"

	"github.com/JuFnd/go-proxy/configs"
	"github.com/JuFnd/go-proxy/internal/app/proxy/pkg/app"
	"github.com/JuFnd/go-proxy/internal/app/proxy/pkg/logger"
	"github.com/JuFnd/go-proxy/internal/app/proxy/server"
	"github.com/JuFnd/go-proxy/internal/app/server/delivery"
	"github.com/JuFnd/go-proxy/internal/app/server/repository"
	"github.com/JuFnd/go-proxy/internal/app/server/usecase"
)

var loggerSingleton logger.Singleton

func main() {
	logger := loggerSingleton.GetLogger()
	app := app.Init()

	srvCfg := configs.GetHTTPSrvConfig(app.ConfigPath)
	tlsCfg := configs.GetTlsConfig(app.ConfigPath)
	apiCfg := configs.GetWebSrvConfig(app.ConfigPsx)

	requestRepo, _ := repository.GetUserRepo(&apiCfg, logger)
	requestUseCase := usecase.NewProxyUseCase(requestRepo)

	proxy := server.New(&srvCfg, &tlsCfg, &apiCfg, requestUseCase, logger)
	api := delivery.GetApi(requestUseCase, proxy, &srvCfg, logger)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := proxy.ListenAndServe(); err != nil {
			logger.Fatalln(err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		if err := api.ListenAndServe(); err != nil {
			logger.Fatalln(err.Error())
		}
	}()

	wg.Wait()
}
