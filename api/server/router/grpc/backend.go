package grpc // import "github.com/rumpl/bof/api/server/router/grpc"

import "google.golang.org/grpc"

// Backend abstracts a registerable GRPC service.
type Backend interface {
	RegisterGRPC(*grpc.Server)
}
