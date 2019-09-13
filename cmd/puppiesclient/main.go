package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	rpc "github.com/moz5691/lambda-proxy/rpc/puppies"
)

func main() {

	var apigatewayURL string

	var client rpc.Puppies
	var weightFloat float64
	var ageFloat float64

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apigatewayURL = os.Getenv("APIGATEWAY_URL")

	if apigatewayURL == "" {
		client = rpc.NewPuppiesProtobufClient("http://localhost:8080", &http.Client{})
	} else {
		client = rpc.NewPuppiesProtobufClient(apigatewayURL, &http.Client{})
	}

	method := flag.String("m", "r", "r-read, d-delete, u-update, s-scan")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	switch *method {
	case "r":
		fmt.Println("=== Read : get a dob by name and id ===")
		fmt.Print("Id: ")
		id, _ := reader.ReadString('\n')
		id = strings.Replace(id, "\n", "", -1)

		fmt.Print("Name: ")
		name, _ := reader.ReadString('\n')
		name = strings.Replace(name, "\n", "", -1)

		term := &rpc.Name{
			Id:   id,
			Name: name,
		}
		res, err := client.GetByName(context.Background(), term)

		if err != nil {
			fmt.Printf("oh no !!!: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("I have found it. : %+v\n", res)

	case "d":
		fmt.Println("=== Delete : delete a dog by name and id ===")
		fmt.Print("Id: ")
		id, _ := reader.ReadString('\n')
		id = strings.Replace(id, "\n", "", -1)

		fmt.Print("Name: ")
		name, _ := reader.ReadString('\n')
		name = strings.Replace(name, "\n", "", -1)

		term := &rpc.Name{
			Id:   id,
			Name: name,
		}

		res, err := client.DeleteByName(context.Background(), term)
		if err != nil {
			fmt.Printf("oh no !!!: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deleted item : %+v\n", res)

	case "s":
		fmt.Println("=== Scan : Scan dogs by breed ===")
		fmt.Print("Breed: ")
		breed, _ := reader.ReadString('\n')
		breed = strings.Replace(breed, "\n", "", -1)

		term := &rpc.Breed{
			Breed: breed,
		}
		res, err := client.ScanBreed(context.Background(), term)

		if err != nil {
			fmt.Printf("oh no.... %+v\n", err)
			os.Exit(1)
		}
		fmt.Printf("scanned items: %+v\n", res)

	case "u":
		fmt.Println("=== Update : update age and weight with a dog's id ==")
		fmt.Print("Id: ")
		id, _ := reader.ReadString('\n')
		id = strings.Replace(id, "\n", "", -1)

		fmt.Print("Name: ")
		name, _ := reader.ReadString('\n')
		name = strings.Replace(name, "\n", "", -1)

		fmt.Print("Age: ")
		age, _ := reader.ReadString('\n')
		age = strings.Replace(age, "\n", "", -1)
		if ageFloat, err = strconv.ParseFloat(age, 64); err != nil {
			fmt.Println("fail to convert to float")
		}

		fmt.Print("Weight: ")
		weight, _ := reader.ReadString('\n')
		weight = strings.Replace(weight, "\n", "", -1)
		if weightFloat, err = strconv.ParseFloat(weight, 64); err != nil {
			fmt.Println("fail to convert to float")
		}

		term := &rpc.Update{
			Id:     id,
			Name:   name,
			Age:    ageFloat,
			Weight: weightFloat,
		}
		res, err := client.UpdateAgeWeight(context.Background(), term)

		if err != nil {
			fmt.Printf("oh no.... %+v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated : %+v\n", res)

	}

}
