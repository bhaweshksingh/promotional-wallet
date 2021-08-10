package server

import (
	"account/pkg/background/handler"
	"account/pkg/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server interface {
	Start()
}

type appServer struct {
	cfg    config.Config
	lgr    *zap.Logger
	router http.Handler
	bgh    *handler.AccountBackgroundHandler
}

func (s *appServer) Start() {
	server := newHTTPServer(s.cfg.GetHTTPServerConfig(), s.router)

	s.lgr.Sugar().Infof("listening on %s", s.cfg.GetHTTPServerConfig().GetAddress())
	go func() { _ = server.ListenAndServe() }()
	done := make(chan bool)
	go s.startDataRefresher(done)()

	waitForShutdown(server, s.lgr, done)
}

func waitForShutdown(server *http.Server, lgr *zap.Logger, done chan bool) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigCh
	done <- true

	defer func() { _ = lgr.Sync() }()

	err := server.Shutdown(context.Background())
	if err != nil {
		lgr.Error(err.Error())
		return
	}

	lgr.Info("server shutdown successful")
}


func (s *appServer) startDataRefresher(done chan bool) func() {
	return func() {
		tickerDuration := time.Duration(s.cfg.GetDataRefresherConfig().GetTickerIntervalInSec())
		ticker := time.NewTicker(tickerDuration * time.Second)

		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				s.lgr.Sugar().Infof("Refreshing Data...", t)
				err := s.bgh.UpdateAccountExpiredEntries()
				if err != nil {
					s.lgr.Error(fmt.Sprintf("Error occurred while updating account info: %v", err.Error()))
				}
			}
		}
	}
}

func newHTTPServer(cfg config.HTTPServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Handler:      handler,
		Addr:         cfg.GetAddress(),
		WriteTimeout: time.Second * time.Duration(cfg.GetReadTimeout()),
		ReadTimeout:  time.Second * time.Duration(cfg.GetWriteTimeout()),
	}
}

func NewServer(cfg config.Config, lgr *zap.Logger, router http.Handler, bgh *handler.AccountBackgroundHandler) Server {
	return &appServer{
		cfg:    cfg,
		lgr:    lgr,
		router: router,
		bgh:    bgh,
	}
}
