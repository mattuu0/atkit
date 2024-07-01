package sdk

import (
	"auth/sdk/sgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func Main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	client := sgrpc.NewSampleServiceClient(conn)

	response, err := client.GetData(context.Background(), &sgrpc.Message{Body: "送信データ"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Print(response.Body)

	defer conn.Close()
}
