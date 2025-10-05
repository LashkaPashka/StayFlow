package payment

import (
	"context"

	paymentV1 "github.com/LashkaPashka/StayFlow/PaymentService/payment_proto/gen/go/payment"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/lib/converter"
	"github.com/LashkaPashka/StayFlow/PaymentService/server/internal/model"
	"google.golang.org/grpc"
)

type Service interface {
	CreatePayment(ctx context.Context, payment *model.Payment) (paymentID, url string,err error)
}

type serverAPI struct {
	paymentV1.UnimplementedPaymentServiceServer
	service Service
}

func Register(serivce Service, gRPCServer *grpc.Server) {
	paymentV1.RegisterPaymentServiceServer(gRPCServer, &serverAPI{service: serivce})
}

func (s *serverAPI) CreatePayment(
	ctx context.Context,
	in *paymentV1.CreatePaymentRequest,
) (*paymentV1.CreatePaymentResponse, error){
	payment := converter.Convert(in)

	paymentID, url, err := s.service.CreatePayment(ctx, payment)
	if err != nil {
		return &paymentV1.CreatePaymentResponse{}, err
	}

	return &paymentV1.CreatePaymentResponse{
		Url: url,
		PaymentId: paymentID,
		Status: string(payment.Status),
	}, nil
}

func (s *serverAPI) GetPaymentStatus(
	ctx context.Context,
	in *paymentV1.GetPaymentStatusRequest,
) (*paymentV1.GetPaymentStatusReponse, error){

	return &paymentV1.GetPaymentStatusReponse{Status: "PENDING"}, nil
}

func (s *serverAPI) RefundPayment(
	ctx context.Context,
	in *paymentV1.RefundPaymentRequest,
) (*paymentV1.RefundPaymentResponse, error){

	return &paymentV1.RefundPaymentResponse{}, nil
}