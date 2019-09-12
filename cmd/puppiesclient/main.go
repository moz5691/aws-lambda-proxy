package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	rpc "github.com/moz5691/lambda-proxy/rpc/puppies"
)

func main() {

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// apigatewayURL = os.Getenv("APIGATEWAY_URL")
	client := rpc.NewPuppiesProtobufClient("https://i39qtt77t5.execute-api.us-east-1.amazonaws.com/dev/", &http.Client{})
	// client := rpc.NewPuppiesProtobufClient("http://localhost:8080", &http.Client{})

	term := &rpc.Name{
		Id:   "00a41382-3231-4089-b84c-301327e48f3e",
		Name: "Brady",
	}

	res, err := client.GetByName(context.Background(), term)

	if err != nil {
		fmt.Printf("oh no !!!: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("I have found it. : %+v\n", res)

	del := &rpc.Name{
		Id:   "11e14757-c27a-40cf-b54a-f7b0ee268327",
		Name: "Samson",
	}

	res1, err1 := client.DeleteByName(context.Background(), del)
	if err1 != nil {
		fmt.Printf("oh no !!!: %+v\n", err1)
		os.Exit(1)
	}

	fmt.Printf("Deleted item : %+v\n", res1)

}
