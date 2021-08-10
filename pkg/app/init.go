package app

import (
	"account/pkg/accountinfo"
	"account/pkg/background/handler"
	"account/pkg/config"
	"account/pkg/http/router"
	"account/pkg/http/server"
	"account/pkg/ledger"
	"account/pkg/reporters"
	"account/pkg/repository"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
)

func initHTTPServer(configFile string) {
	config := config.NewConfig(configFile)
	logger := initLogger(config)
	rt := initRouter(config, logger)
	bgh := initBackgroundHandler(config, logger)

	server.NewServer(config, logger, rt, bgh).Start()
}

func initRouter(cfg config.Config, logger *zap.Logger) http.Handler {
	accountInfoRepository, ledgerRepository := initRepository(cfg)
	accountInfoService, ledgerService := initService(accountInfoRepository, ledgerRepository)

	return router.NewRouter(logger, accountInfoService, ledgerService)
}

func initService(accountInfoRepository repository.AccountRepository, ledgerRepository repository.LedgerRepository) (accountinfo.Service, ledger.Service) {
	accountInfoService := accountinfo.NewAccountInfoService(accountInfoRepository)
	ledgerService := ledger.NewLedgerService(ledgerRepository)

	return accountInfoService, ledgerService
}

func initBackgroundHandler(cfg config.Config, logger *zap.Logger) *handler.AccountBackgroundHandler {
	accountInfoRepository, ledgerRepository := initRepository(cfg)
	accountInfoService, ledgerService := initService(accountInfoRepository, ledgerRepository)


	return handler.NewAccountBackgroundHandler(logger, accountInfoService, ledgerService, cfg.GetDataRefresherConfig())
}

func initRepository(cfg config.Config) (repository.AccountRepository, repository.LedgerRepository) {
	dbConfig := cfg.GetDBConfig()
	dbHandler := repository.NewDBHandler(dbConfig)

	db, err := dbHandler.GetDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	return repository.NewAccountInfoRepository(db), repository.NewLedgerRepository(db)
}

func initLogger(cfg config.Config) *zap.Logger {
	return reporters.NewLogger(
		cfg.GetLogConfig().GetLevel(),
		getWriters(cfg.GetLogFileConfig())...,
	)
}

func getWriters(cfg config.LogFileConfig) []io.Writer {
	return []io.Writer{
		os.Stdout,
		reporters.NewExternalLogFile(cfg),
	}
}
