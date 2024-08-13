package helper

import (
	"context"
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"google.golang.org/grpc/metadata"
)

func GetMetadata(ctx context.Context) *entity.Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return new(entity.Metadata)
	}

	m := new(entity.Metadata)

	hosts := md.Get("Host")
	if len(hosts) > 0 {
		m.Host = hosts[0]
	}

	ips := md.Get("X-Forwarded-For")
	if len(ips) > 0 {
		m.Ip = ips[0]
	}

	protocols := md.Get("X-Forwarded-Proto")
	if len(protocols) > 0 {
		m.Protocol = protocols[0]
	}

	return m
}
