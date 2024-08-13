package middleware

import "github.com/dwprz/prasorganic-user-service/src/core/restful/client"

type Middleware struct {
	restfulClient *client.Restful
}

func New(rc *client.Restful) *Middleware {
	return &Middleware{
		restfulClient: rc,
	}
}
