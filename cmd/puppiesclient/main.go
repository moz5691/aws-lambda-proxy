package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	rpc "github.com/moz5691/lambda-proxy/rpc/puppies"
)

func main() {
	client := rpc.NewPuppiesProtobufClient("https://i39qtt77t5.execute-api.us-east-1.amazonaws.com/dev/", &http.Client{})
	// client := rpc.NewPuppiesProtobufClient("http://localhost:8080", &http.Client{})

	term := &rpc.GetByIdReq{
		Id:   "00a41382-3231-4089-b84c-301327e48f3e",
		Name: "Brady",
	}

	res, err := client.GetById(context.Background(), term)

	if err != nil {
		fmt.Printf("oh no !!!: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("I have found it. : %+v\n", res)
}
