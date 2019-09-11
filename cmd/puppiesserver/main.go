package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/rs/cors"

	"github.com/moz5691/lambda-proxy/internal/puppiesserver"
	"github.com/moz5691/lambda-proxy/rpc/puppies"
)

func main() {

	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}
	mux := http.NewServeMux()
	server := &puppiesserver.Server{}

	puppiesHandler := puppies.NewPuppiesServer(server, nil)
	mux.Handle(puppies.PuppiesPathPrefix, puppiesHandler)

	listenPort, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		listenPort = "8080"
	}

	handler := cors.AllowAll().Handler(mux)

	awsExecutionEnv, exists := os.LookupEnv("AWS_EXECUTION_ENV")

	if !exists {
		fmt.Println("Locally on port", listenPort)
		http.ListenAndServe(":8080", handler)
	} else {
		fmt.Println(awsExecutionEnv)
		log.Fatal(gateway.ListenAndServe("", mux))
	}
}
