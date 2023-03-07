package grpc

import (
	"context"
	"io"
	"log"
	"time"

	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
	"github.com/timpamungkas/grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"
)

func (a *GrpcAdapter) GetCurrentBalance(
	ctx context.Context, in *bank.CurrentBalanceRequest) (*bank.CurrentBalanceResponse, error) {
	now := time.Now()
	bal := a.bankService.FindCurrentBalance(in.AccountNumber)

	return &bank.CurrentBalanceResponse{
		Amount: bal,
		CurrentDate: &date.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(now.Day()),
		},
	}, nil
}

func (a *GrpcAdapter) FetchExchangeRates(in *bank.ExchangeRateRequest,
	stream bank.BankService_FetchExchangeRatesServer) error {
	for {
		rate := a.bankService.FindExchangeRate(in.FromCurrency, in.ToCurrency, time.Now())

		stream.Send(
			&bank.ExchangeRateResponse{
				FromCurrency: in.FromCurrency,
				ToCurrency:   in.ToCurrency,
				Rate:         rate,
				Timestamp:    time.Now().Format(time.RFC3339),
			},
		)

		time.Sleep(3 * time.Second)
	}
}

func (a *GrpcAdapter) SummarizeTransactions(stream bank.BankService_SummarizeTransactionsServer) error {
	tsum := dbank.TransactionSummary{
		SummaryOnDate: time.Now(),
		SumIn:         0,
		SumOut:        0,
		SumTotal:      0,
	}
	acct := ""

	for {
		in, err := stream.Recv()

		if err == io.EOF {
			res := bank.TransactionSummary{
				AccountNumber: acct,
				SumAmountIn:   tsum.SumIn,
				SumAmountOut:  tsum.SumOut,
				SumTotal:      tsum.SumTotal,
				TransactionDate: &date.Date{
					Year:  int32(tsum.SummaryOnDate.Year()),
					Month: int32(tsum.SummaryOnDate.Month()),
					Day:   int32(tsum.SummaryOnDate.Day()),
				},
			}

			return stream.SendAndClose(
				&res,
			)
		}

		if err != nil {
			log.Fatalf("Error while reading from client : %v", err)
		}

		acct = in.AccountNumber
		ts, err := toTime(in.Timestamp)

		if err != nil {
			log.Fatalf("Error while parsing timestamp %v : %v", in.Timestamp, err)
		}

		ttype := dbank.TransactionStatusUnknown

		if in.Type == bank.TransactionType_TRANSACTION_TYPE_IN {
			ttype = dbank.TransactionStatusIn
		} else if in.Type == bank.TransactionType_TRANSACTION_TYPE_OUT {
			ttype = dbank.TransactionStatusOut
		}

		tcur := dbank.Transaction{
			Amount:          in.Amount,
			Timestamp:       ts,
			TransactionType: ttype,
		}

		err = a.bankService.CalculateTransactionSummary(&tsum, tcur)

		if err != nil {
			return err
		}
	}
}

func toTime(dt *datetime.DateTime) (time.Time, error) {
	res := time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day),
		int(dt.Hours), int(dt.Minutes), int(dt.Seconds), int(dt.Nanos), time.UTC)

	return res, nil
}

func (a *GrpcAdapter) TransferMultiple(stream bank.BankService_TransferMultipleServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading from client : %v", err)
		}

		tt := dbank.TransferTransaction{
			FromAccountNumber: req.FromAccountNumber,
			ToAccountNumber:   req.ToAccountNumber,
			Currency:          req.Currency,
			Amount:            req.Amount,
		}

		_, transferSuccess, err := a.bankService.Transfer(tt)

		if err != nil {
			return err
		}

		res := bank.TransferResponse{
			FromAccountNumber: req.FromAccountNumber,
			ToAccountNumber:   req.ToAccountNumber,
			Amount:            req.Amount,
			Timestamp:         currentDatetime(),
		}

		if transferSuccess {
			res.Status = bank.TransferStatus_TRANSFER_STATUS_SUCCESS
		} else {
			res.Status = bank.TransferStatus_TRANSFER_STATUS_FAILED
		}

		err = stream.Send(&res)

		if err != nil {
			log.Fatalf("Error while sending response to client : %v", err)
		}
	}

}

func currentDatetime() *datetime.DateTime {
	now := time.Now()

	return &datetime.DateTime{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		Day:        int32(now.Day()),
		Hours:      int32(now.Hour()),
		Minutes:    int32(now.Minute()),
		Seconds:    int32(now.Second()),
		Nanos:      int32(now.Second()),
		TimeOffset: &datetime.DateTime_UtcOffset{},
	}
}
