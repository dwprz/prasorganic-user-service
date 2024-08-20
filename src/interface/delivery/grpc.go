package delivery

import (
	"context"
	pb "github.com/dwprz/prasorganic-proto/protogen/otp"
)

type OtpGrpc interface {
	Send(ctx context.Context, email string) error
	Verify(ctx context.Context, data *pb.VerifyReq) (*pb.VerifyRes, error)
}
