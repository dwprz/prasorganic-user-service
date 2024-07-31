package helper

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type Metadata struct {
	Host     string
	Ip       string
	Protocol string
}

func GetMetadata(ctx context.Context) *Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return new(Metadata)
	}

	m := new(Metadata)

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
