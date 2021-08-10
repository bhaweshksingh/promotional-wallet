package ledger

import (
	"account/pkg/accountinfo/dto"
	"account/pkg/ledger/model"
	"account/pkg/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"sort"
	"time"
)

type Service interface {
	CreateLedgerEntry(ctx context.Context, info *model.Ledger) error
	GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error)
	ExpireCredits(ctx context.Context, accountID uuid.UUID) ([]*model.Ledger, error)
	AddDebitEntry(ctx context.Context, debitEntry *model.Ledger) error
}

type ledgerService struct {
	repository repository.LedgerRepository
}

func (ls *ledgerService) CreateLedgerEntry(ctx context.Context, info *model.Ledger) error {
	if info.Activity == dto.Credit {
		err := ls.repository.CreateLedgerEntries(ctx, []*model.Ledger{info})
		if err != nil {
			return fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
		}
	} else if info.Activity == dto.Debit {

	}

	return nil
}

func (ls *ledgerService) GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error) {
	entries, err := ls.repository.GetEntries(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
	}

	return entries, nil
}

// any priority level that has the credit and some value remaining from the debit and expiry
// and has crossed the timing is eligible to be  expired
func (ls *ledgerService) ExpireCredits(ctx context.Context, accountID uuid.UUID) ([]*model.Ledger, error) {
	_, totalCredit, err := ls.fetchAggregateEntries(ctx, accountID.String())
	if err != nil {
		return nil, err
	}
	expiryEntries := make([]*model.Ledger, 0)
	for priority, credit := range totalCredit {
		if credit.Expiry.Before(time.Now()) && credit.Amount > 0 {
			var expiryEntry *model.Ledger
			expiryEntry = &model.Ledger{
				AccountID: credit.AccountID,
				Amount:    credit.Amount,
				Priority:  priority,
				Activity:  dto.Expiration,
				Expiry:    time.Now(),
				CreatedAt: time.Now(),
			}
			expiryEntries = append(expiryEntries, expiryEntry)
		}

	}
	return expiryEntries, nil
}
func (ls *ledgerService) AddDebitEntry(ctx context.Context, debitEntryRequest *model.Ledger) error {
	_, totalCredit, err := ls.fetchAggregateEntries(ctx, debitEntryRequest.AccountID.String())
	if err != nil {
		return err
	}
	targetDebit := debitEntryRequest.Amount
	debitEntries := make([]*model.Ledger, 0)

	for priority, credit := range totalCredit {
		if targetDebit <= 0 {
			break
		}
		var debitEntry *model.Ledger
		if credit.Amount <= targetDebit {
			debitEntry = &model.Ledger{
				AccountID: debitEntryRequest.AccountID,
				Amount:    credit.Amount,
				Priority:  priority,
				Activity:  dto.Debit,
				Expiry:    time.Now(),
				CreatedAt: time.Now(),
			}
			targetDebit -= credit.Amount
		} else {
			debitEntry = &model.Ledger{
				AccountID: debitEntryRequest.AccountID,
				Amount:    credit.Amount - targetDebit,
				Priority:  priority,
				Activity:  dto.Debit,
				Expiry:    time.Now(),
				CreatedAt: time.Now(),
			}
			targetDebit = 0
		}
		debitEntries = append(debitEntries, debitEntry)
	}
	if targetDebit != 0 {
		return fmt.Errorf("Service.DebitRequest failed. Not enough credits. Error: %w", err)
	}
	err = ls.repository.CreateLedgerEntries(ctx, debitEntries)
	if err != nil {
		return fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
	}

	return nil
}

type ExpirableAmount struct {
	Amount    int64
	Expiry    time.Time
	AccountID uuid.UUID
}

func (ls *ledgerService) fetchAggregateEntries(ctx context.Context, accountID string) (map[int64]map[string]*model.AggregateEntry, map[int64]*ExpirableAmount, error) {
	aggregateEntries, err := ls.repository.GetEntriesByPriority(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}
	groupedAggregateEntries := groupByPriorityAndType(aggregateEntries)
	sortedPriorities := getSortedKeys(groupedAggregateEntries)

	totalCredit := make(map[int64]*ExpirableAmount, len(sortedPriorities))
	for _, priority := range sortedPriorities {
		totalCredit[priority] = &ExpirableAmount{}
	}
	// read the debit, credit and expired
	// credit minus expired grouped by priority
	for _, priority := range sortedPriorities {
		// Assumption is that no negative credits are there per priority.
		// The current logic is ensuring that.
		if priorityEntries, ok := groupedAggregateEntries[priority]; ok {
			if creditEntry, ok := priorityEntries[dto.Credit]; ok {
				totalCredit[creditEntry.Priority].Amount += creditEntry.Amount
				totalCredit[creditEntry.Priority].Expiry = creditEntry.Expiry
				totalCredit[creditEntry.Priority].AccountID = creditEntry.AccountID
			}
			if debitEntry, ok := priorityEntries[dto.Debit]; ok {
				totalCredit[debitEntry.Priority].Amount -= debitEntry.Amount
				totalCredit[debitEntry.Priority].Expiry = debitEntry.Expiry
				totalCredit[debitEntry.Priority].AccountID = debitEntry.AccountID
			}
			if expiredEntry, ok := priorityEntries[dto.Expiration]; ok {
				totalCredit[expiredEntry.Priority].Amount -= expiredEntry.Amount
				totalCredit[expiredEntry.Priority].Expiry = expiredEntry.Expiry
				totalCredit[expiredEntry.Priority].AccountID = expiredEntry.AccountID
			}
		}
	}
	return groupedAggregateEntries, totalCredit, err
}

func getSortedKeys(entries map[int64]map[string]*model.AggregateEntry) []int64 {
	priorities := make([]int64, 0)
	for priority, _ := range entries {
		priorities = append(priorities, priority)
	}
	int64AsIntValues := make([]int, len(priorities))

	for i, val := range priorities {
		int64AsIntValues[i] = int(val)
	}
	sort.Ints(int64AsIntValues)

	for i, val := range int64AsIntValues {
		priorities[i] = int64(val)
	}
	return priorities
}

func groupByPriorityAndType(entries []*model.AggregateEntry) map[int64]map[string]*model.AggregateEntry {
	groupedEntries := make(map[int64]map[string]*model.AggregateEntry, 0)
	for _, entry := range entries {
		if prioritisedEntries, ok := groupedEntries[entry.Priority]; !ok {
			newPrioritisedEntries := make(map[string]*model.AggregateEntry, 0)
			newPrioritisedEntries[entry.Activity] = entry
			groupedEntries[entry.Priority] = newPrioritisedEntries
		} else {
			prioritisedEntries[entry.Activity] = entry
		}
	}

	return groupedEntries
}
func NewLedgerService(repository repository.LedgerRepository) Service {
	return &ledgerService{
		repository: repository,
	}
}
