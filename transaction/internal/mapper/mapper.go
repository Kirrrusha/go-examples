package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"transaction/internal/model"
	transactionpb "transaction/pkg/transaction/go"
)

func TransactionToPb(transaction model.Transaction) *transactionpb.Transaction {
	return &transactionpb.Transaction{
		Id:        transaction.ID,
		UserId:    transaction.UserID,
		Amount:    transaction.Amount,
		Type:      string(transaction.Type),
		Status:    string(transaction.Status),
		CreatedAt: timestamppb.New(transaction.CreatedAt),
		UpdatedAt: timestamppb.New(transaction.UpdatedAt),
	}
}

func TransactionEntryToPb(entry model.TransactionEntry) *transactionpb.TransactionEntry {
	return &transactionpb.TransactionEntry{
		Id:            entry.ID,
		TransactionId: entry.TransactionID,
		AccountId:     entry.AccountID,
		Direction:     string(entry.Direction),
		Amount:        entry.Amount,
		CreatedAt:     timestamppb.New(entry.CreatedAt),
		UpdatedAt:     timestamppb.New(entry.UpdatedAt),
	}
}

func TransactionDetailsToPb(details model.TransactionDetails) *transactionpb.TransactionDetails {
	entries := make([]*transactionpb.TransactionEntry, len(details.Entries))
	for i, entry := range details.Entries {
		entries[i] = TransactionEntryToPb(entry)
	}

	return &transactionpb.TransactionDetails{
		Transaction: TransactionToPb(details.Transaction),
		Entries:     entries,
	}
}

func TransactionDetailsListToPb(details []model.TransactionDetails) []*transactionpb.TransactionDetails {
	res := make([]*transactionpb.TransactionDetails, len(details))
	for i, item := range details {
		res[i] = TransactionDetailsToPb(item)
	}

	return res
}

func GetTransactionsRequestToParams(req *transactionpb.GetTransactionsRequest) model.GetTransactionsParams {
	params := model.GetTransactionsParams{
		Limit:  100,
		Offset: 0,
	}
	if req == nil {
		return params
	}

	if req.UserId != nil {
		userID := req.GetUserId()
		params.UserID = &userID
	}
	if req.Type != nil {
		txType := req.GetType()
		params.Type = &txType
	}
	if req.Status != nil {
		status := req.GetStatus()
		params.Status = &status
	}
	if req.GetDateFrom() != nil {
		dateFrom := req.GetDateFrom().AsTime()
		params.DateFrom = &dateFrom
	}
	if req.GetDateTo() != nil {
		dateTo := req.GetDateTo().AsTime()
		params.DateTo = &dateTo
	}
	if req.GetPagination() != nil {
		params.Limit = int(req.GetPagination().GetLimit())
		params.Offset = int(req.GetPagination().GetOffset())
	}

	return params
}
