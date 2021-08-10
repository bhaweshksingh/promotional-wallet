package accountinfo

import (
	"account/pkg/accountinfo/dto"
	"account/pkg/accountinfo/model"
	"account/pkg/repository"
	"context"
	"fmt"
)

type Service interface {
	CreateOrUpdateAccountInfo(ctx context.Context, info *model.AccountInfo) error
	GetAccountsFor(ctx context.Context, accountQuery *dto.AccountQuery) ([]model.AccountInfo, error)
}

type accountInfoService struct {
	repository repository.AccountRepository
}

func (ais *accountInfoService) CreateOrUpdateAccountInfo(ctx context.Context, accountInfoEvent *model.AccountInfo) error {
	accountQuery := &dto.AccountQuery{
		UserID:    accountInfoEvent.UserID,
		AccountID: accountInfoEvent.ID,
	}
	accountInfos, err := ais.repository.GetAccountData(ctx, accountQuery)
	fmt.Println(accountQuery)
	if err != nil {
		return fmt.Errorf("Service.GetAccountsFor", err)
	}
	if len(accountInfos) < 1 {
		return fmt.Errorf("Service.GetAccountsFor: No Account exists for the userid", err)
	}
	if len(accountInfos) > 1 {
		return fmt.Errorf("Service.GetAccountsFor: User should not have 2 accounts.", err)
	}
	existingAccountInfo := accountInfos[0]
	existingAccountInfo.UpdateBalance(accountInfoEvent.Balance)
	err = ais.repository.CreateOrUpdateAccountInfo(ctx, &existingAccountInfo)
	if err != nil {
		return fmt.Errorf("Service.CreateOrUpdateAccountInfo failed. Error: %w", err)
	}
	return nil
}

func (ais *accountInfoService) GetAccountsFor(ctx context.Context, accountQuery *dto.AccountQuery) ([]model.AccountInfo, error) {
	fmt.Println(accountQuery)
	accountInfos, err := ais.repository.GetAccountData(ctx, accountQuery)
	if err != nil {
		return nil, fmt.Errorf("Service.GetAccountsFor", err)
	}

	return accountInfos, nil
}

func NewAccountInfoService(repository repository.AccountRepository) Service {
	return &accountInfoService{
		repository: repository,
	}
}
