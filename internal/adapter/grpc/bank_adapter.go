package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
	"github.com/timpamungkas/grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *GrpcAdapter) GetCurrentBalance(
	ctx context.Context, in *bank.CurrentBalanceRequest) (*bank.CurrentBalanceResponse, error) {
	now := time.Now()
	bal, err := a.bankService.FindCurrentBalance(in.AccountNumber)

	if err != nil {
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"account %v not found", in.AccountNumber,
		)
	}

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
		now := time.Now().Truncate(time.Second)
		rate, err := a.bankService.FindExchangeRate(in.FromCurrency, in.ToCurrency, now)

		if err != nil {
			s := status.New(codes.InvalidArgument,
				"Currency not valid. Please pass valid currency for both from and to")
			s, _ = s.WithDetails(&errdetails.ErrorInfo{
				Domain: "my-bank-website.com",
				Reason: "INVALID_CURRENCY",
				Metadata: map[string]string{
					"from_currency": in.FromCurrency,
					"to_currency":   in.ToCurrency,
				},
			})

			return s.Err()
		}

		stream.Send(
			&bank.ExchangeRateResponse{
				FromCurrency: in.FromCurrency,
				ToCurrency:   in.ToCurrency,
				Rate:         rate,
				Timestamp:    now.Format(time.RFC3339),
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

		accountUuid, err := a.bankService.CreateTransaction(in.AccountNumber, tcur)

		if err != nil && accountUuid == uuid.Nil {
			s := status.New(codes.InvalidArgument, err.Error())
			s, _ = s.WithDetails(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "account_number",
						Description: "Invalid account number",
					},
				},
			})

			return s.Err()
		} else if err != nil && accountUuid != uuid.Nil {
			s := status.New(codes.InvalidArgument, err.Error())
			s, _ = s.WithDetails(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field: "amount",
						Description: fmt.Sprintf(
							"Requested amount %v exceed available balance", in.Amount),
					},
				},
			})

			return s.Err()
		}

		err = a.bankService.CalculateTransactionSummary(&tsum, tcur)

		if err != nil {
			return err
		}
	}
}

func toTime(dt *datetime.DateTime) (time.Time, error) {
	if dt == nil {
		now := time.Now()

		dt = &datetime.DateTime{
			Year:    int32(now.Year()),
			Month:   int32(now.Month()),
			Day:     int32(now.Day()),
			Hours:   int32(now.Hour()),
			Minutes: int32(now.Minute()),
			Seconds: int32(now.Second()),
			Nanos:   int32(now.Nanosecond()),
		}
	}

	res := time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day),
		int(dt.Hours), int(dt.Minutes), int(dt.Seconds), int(dt.Nanos), time.UTC)

	return res, nil
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
			Currency:          req.Currency,
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
