package handler

import (
	"account/pkg/accountinfo"
	"account/pkg/accountinfo/dto"
	"account/pkg/http/contract"
	"account/pkg/http/internal/resperr"
	"account/pkg/http/internal/utils"
	"account/pkg/ledger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"

	"go.uber.org/zap"
)

type AccountHandler struct {
	lgr        *zap.Logger
	accountSvc accountinfo.Service
	ledgerSvc  ledger.Service
}

func NewAccountHandler(lgr *zap.Logger, accountSvc accountinfo.Service, ledgerSvc ledger.Service) *AccountHandler {
	return &AccountHandler{
		lgr:        lgr,
		accountSvc: accountSvc,
		ledgerSvc:  ledgerSvc,
	}
}

func (ah *AccountHandler) Credit(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	var ae dto.AccountEvent
	err := utils.ParseRequest(req, &ae)
	if err != nil {
		return err
	}
	userID, _ := uuid.Parse(ae.UserID)
	accounts, err := ah.accountSvc.GetAccountsFor(ctx, &dto.AccountQuery{UserID: userID})
	if err != nil {
		return err
	}
	err = ah.ledgerSvc.CreateLedgerEntry(ctx, ae.GetCreditLedgerEntry(accounts[0]))
	if err != nil {
		return fmt.Errorf("AccountHandler.CreateLedgerEntry . error %v", err)
	}
	err = ah.accountSvc.CreateOrUpdateAccountInfo(ctx, ae.GetAccountInfo(accounts[0], dto.Credit))
	if err != nil {
		return fmt.Errorf("AccountHandler.CreateLedgerEntry . error %v", err)
	}

	ah.lgr.Debug("msg", zap.String("eventCode", utils.AccountInfoUpdated))
	utils.WriteSuccessResponse(resp, http.StatusCreated, contract.AccountInfoCreationSuccess)
	return nil
}

func (ah *AccountHandler) GetLogs(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	params := mux.Vars(req)
	fmt.Println(params)
	activityType := params["activity"]
	if ok, _ := dto.GetAllowedActivityTypes()[activityType]; activityType != "" && !ok {
		utils.WriteFailureResponse(resp, resperr.NewResponseError(http.StatusBadRequest, contract.BadLogsRequest))
		return nil
	}
	logQuery := &dto.LogQuery{
		ActivityType: activityType,
	}
	fmt.Println(logQuery)

	entries, err := ah.ledgerSvc.GetEntries(ctx, logQuery)
	if err != nil {
		return err
	}

	ah.lgr.Debug("msg", zap.String("eventCode", utils.LedgerEntriesFetched))
	utils.WriteSuccessResponse(resp, http.StatusOK, entries)
	return nil
}

func (ah *AccountHandler) Debit(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	var ae dto.AccountEvent
	err := utils.ParseRequest(req, &ae)
	if err != nil {
		return err
	}
	userID, _ := uuid.Parse(ae.UserID)
	accounts, err := ah.accountSvc.GetAccountsFor(ctx, &dto.AccountQuery{UserID: userID})
	if err != nil {
		return err
	}
	existingAccount := accounts[0]
	if existingAccount.Balance < ae.Amount {
		utils.WriteFailureResponse(resp, resperr.NewResponseError(http.StatusPreconditionFailed, contract.AccountDoesntHaveEnoughBalance))
	}
	err = ah.ledgerSvc.AddDebitEntry(ctx, ae.GetDebitLedgerEntry(accounts[0]))
	if err != nil {
		return fmt.Errorf("AccountHandler.CreateLedgerEntry . error %v", err)
	}
	err = ah.accountSvc.CreateOrUpdateAccountInfo(ctx, ae.GetAccountInfo(accounts[0], dto.Debit))
	if err != nil {
		return fmt.Errorf("AccountHandler.CreateLedgerEntry . error %v", err)
	}

	ah.lgr.Debug("msg", zap.String("eventCode", utils.AccountInfoUpdated))
	utils.WriteSuccessResponse(resp, http.StatusCreated, contract.AccountInfoCreationSuccess)
	return nil
}

func (ah *AccountHandler) GetAccountInfoFor(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	params := mux.Vars(req)
	fmt.Println(params)
	userID, _ := uuid.Parse(params["userid"])
	accounts, err := ah.accountSvc.GetAccountsFor(ctx, &dto.AccountQuery{UserID: userID})
	if err != nil {
		return err
	}
	existingAccount := accounts[0]

	ah.lgr.Debug("msg", zap.String("eventCode", utils.AccountInfoUpdated))
	utils.WriteSuccessResponse(resp, http.StatusOK, existingAccount)
	return nil
}
