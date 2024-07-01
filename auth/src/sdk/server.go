package sdk

import (
	"auth/sdk/sgrpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
)

func Init() {

}

func Start_server() {
	//サーバー起動
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT")))

	//エラー処理
	if err != nil {
		panic(err)
	}

	//サーバー種痘
	grpc_server := grpc.NewServer()

	//サーバー登録
	sgrpc.RegisterSampleServiceServer(grpc_server, &Sdk_Server{})

	//サーバー起動
	if err := grpc_server.Serve(listen); err != nil {
		panic(err)
	}
}

type Sdk_Server struct{}

// GetData implements sgrpc.SampleServiceServer.
func (server *Sdk_Server) GetData(context.Context, *sgrpc.Message) (*sgrpc.Message, error) {
	return &sgrpc.Message{Body: "Hello World"}, nil
}
