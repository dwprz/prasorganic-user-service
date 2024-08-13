package interceptor

import (
	"context"

	"github.com/dwprz/prasorganic-user-service/src/common/errors"
	"github.com/dwprz/prasorganic-user-service/src/common/helper"
	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*

Flow kerja server interceptor:
	- semua code di dalam interceptor yang sebelum function handler(ctx, req) dipanggil,
		akan di eksekusi sebelum rpc handler nya dieksekusi.
	- semua code di dalam interceptor yang sesudah function handler(ctx, req) dipanggil,
		akan di eksekusi setelah rpc handler nya dieksekusi.
	- urutannya, jika ada interceptor A, interceptor B dan interceptor C :
		- sebelum rpc handler dieksekusi interceptor A => interceptor B => interceptor C
		- setelah rpc handler dieksekusi interceptor C => interceptor B => interceptor A

*/

type UnaryResponse struct{}

func NewUnaryResponse() *UnaryResponse {
	return &UnaryResponse{}
}

func (u *UnaryResponse) Error(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	res, err := handler(ctx, req)

	if err != nil {
		m := helper.GetMetadata(ctx)

		log.Logger.WithFields(logrus.Fields{
			"host":     m.Host,
			"ip":       m.Ip,
			"protocol": m.Protocol,
			"location": info.FullMethod,
			"from":     "Error interceptor",
		}).Error(err)

		// validation error handling
		if errVldtn, ok := err.(validator.ValidationErrors); ok {
			s := status.New(codes.InvalidArgument, err.Error())

			s, _ = s.WithDetails(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       errVldtn[0].Field(),
						Description: errVldtn[0].Error(),
					},
				},
			})

			return nil, s.Err()
		}

		if errRspn, ok := err.(*errors.Response); ok {
			return nil, status.Errorf(errRspn.GrpcCode, errRspn.Message)
		}

		return nil, status.Errorf(codes.Internal, "sorry, internal server error try again later")
	}

	return res, nil
}

func (u *UnaryResponse) Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			m := helper.GetMetadata(ctx)

			log.Logger.WithFields(logrus.Fields{
				"host":     m.Host,
				"ip":       m.Ip,
				"protocol": m.Protocol,
				"location": info.FullMethod,
				"from":     "Recovery interceptor",
			}).Error(r)

			resp = nil
			err = status.Error(codes.Internal, "sorry, internal server error try again later")
		}
	}()

	res, err := handler(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
