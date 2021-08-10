package router

import (
	 "account/pkg/accountinfo"
	"account/pkg/http/internal/handler"
	"account/pkg/http/internal/middleware"
	"account/pkg/ledger"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	accountCreditPath = "/account/credit"
	accountLogsPath = "/account/logs"
	accountPath = "/account/{userid}"
	accountDebitPath = "/account/debit"
)

func NewRouter(lgr *zap.Logger, accountInfoService accountinfo.Service, ledgerService ledger.Service)  http.Handler {
	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler())

	ah := handler.NewAccountHandler(lgr, accountInfoService, ledgerService)

	router.HandleFunc(accountCreditPath, withMiddlewares(lgr, middleware.WithErrorHandler(lgr, ah.Credit))).Methods(http.MethodPost)
	router.HandleFunc(accountDebitPath, withMiddlewares(lgr, middleware.WithErrorHandler(lgr, ah.Debit))).Methods(http.MethodPost)
	router.HandleFunc(accountLogsPath, withMiddlewares(lgr, middleware.WithErrorHandler(lgr, ah.GetLogs))).
		Methods(http.MethodGet).
		Queries("activity", "{activity}")
	router.HandleFunc(accountPath, withMiddlewares(lgr, middleware.WithErrorHandler(lgr, ah.GetAccountInfoFor))).
		Methods(http.MethodGet)
	return router
}

func withMiddlewares(lgr *zap.Logger, hnd http.HandlerFunc) http.HandlerFunc {
	return middleware.WithSecurityHeaders(middleware.WithReqResLog(lgr, hnd))
}
