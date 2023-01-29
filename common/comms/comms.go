package comms

import (
	"fmt"

	userpb "github.com/iamvasanth07/showcase/common/protos/user"
	videopb "github.com/iamvasanth07/showcase/common/protos/video"
	"google.golang.org/grpc"
)

type gRPCSettings struct {
	GRPC_SERVICE_NAME string
	GRPC_HOST         string
	GRPC_PORT         string
}

func NewGrpcSettings(service string, host string, port string) *gRPCSettings {
	return &gRPCSettings{
		GRPC_SERVICE_NAME: service,
		GRPC_HOST:         host,
		GRPC_PORT:         port,
	}
}

func (g *gRPCSettings) CreateGRPCConn() *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", g.GRPC_HOST, g.GRPC_PORT), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

func (g *gRPCSettings) GetGRPCClient() interface{} {
	conn := g.CreateGRPCConn()
	switch g.GRPC_SERVICE_NAME {
	case "user":
		return userpb.NewUserServiceClient(conn)
	case "video":
		return videopb.NewVideoServiceClient(conn)
	}
	return nil
}
