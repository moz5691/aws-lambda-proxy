package puppiesserver

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	// "math/rand"

	"github.com/moz5691/aws-lambda-proxy/internal/awsutils"
	rpc "github.com/moz5691/aws-lambda-proxy/rpc/puppies"
	"github.com/twitchtv/twirp"
)

type puppy struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Age            float64 `json:"age"`
	Weight         float64 `json:"weight"`
	PrimaryColor   string  `json:"primaryColor"`
	SecondaryColor string  `json:"secondaryColor"`
	Owner          string  `json:"owner"`
	Location       string  `json:"location"`
	Motto          string  `json:"motto"`
	Breed          string  `json:"breed"`
}

type Server struct{}

func (s *Server) GetByName(ctx context.Context, req *rpc.Name) (*rpc.Puppy, error) {
	p := &puppy{}

	fmt.Println("req: %v\n", req)
	av, err := awsutils.MarshalMap(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("av: %v\n", av)

	input := &dynamodb.GetItemInput{
		Key:       av,
		TableName: &awsutils.DynamoTableName,
	}

	res, err := awsutils.GetItem(input)

	fmt.Printf("res: %+v\n", res)

	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return nil, twirp.NewError(twirp.NotFound, "Error")
	}

	if len(res.Item) == 0 {
		return nil, twirp.NewError(twirp.NotFound, "No item found.")
	}

	err = awsutils.UnmarshalMap(res.Item, p)
	if err != nil {
		return nil, err
	}

	return &rpc.Puppy{
		Id:             p.ID,
		Name:           p.Name,
		Age:            p.Age,
		Weight:         p.Weight,
		PrimaryColor:   p.PrimaryColor,
		SecondaryColor: p.SecondaryColor,
		Owner:          p.Owner,
		Location:       p.Location,
		Motto:          p.Motto,
		Breed:          p.Breed,
	}, nil

}

func (s *Server) DeleteByName(ctx context.Context, req *rpc.Name) (*rpc.Name, error) {

	fmt.Printf("req: %v\n", req)
	av, err := awsutils.MarshalMap(req)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: &awsutils.DynamoTableName,
	}

	res, err := awsutils.DeleteItem(input)

	fmt.Printf("delete res: %v\n", res)

	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return nil, twirp.NewError(twirp.NotFound, "Error")
	}

	return &rpc.Name{
		Id:   req.Id,
		Name: req.Name,
	}, nil
}

func (s *Server) UpdateAgeWeight(ctx context.Context, req *rpc.Update) (*rpc.Update, error) {
	fmt.Printf("req: %v\n", req)
	id := req.Id
	name := req.Name
	age := req.Age
	weight := req.Weight

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				N: aws.String(strconv.FormatFloat(age, 'E', -1, 64)),
			},
			":w": {
				N: aws.String(strconv.FormatFloat(weight, 'E', -1, 64)),
			},
		},
		TableName: &awsutils.DynamoTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"name": {
				S: aws.String(name),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set age = :a, weight = :w "),
	}

	_, err := awsutils.UpdateItem(input)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, twirp.NewError(twirp.NotFound, "Error")
	}

	return &rpc.Update{
		Id:     id,
		Name:   name,
		Weight: weight,
		Age:    age,
	}, nil

}

func (s *Server) ScanBreed(ctx context.Context, req *rpc.Breed) (*rpc.PuppyList, error) {
	// maxAge := req.GetAge()
	breed := req.GetBreed()

	filt := expression.Name("breed").Equal(expression.Value(breed))
	// condition := expression.And(expression.Name("breed").Equal(expression.Value(breed)),
	// 	expression.Name("age").LessThanEqual(expression.Value(maxAge)))
	proj := expression.NamesList(expression.Name("name"), expression.Name("age"), expression.Name("breed"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &awsutils.DynamoTableName,
	}

	result, err := awsutils.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("scan result: %v\n", result)
	// numItems := 0

	// for _, i := range result.Items {
	// 	p := &puppy{}

	// 	err = awsutils.UnmarshalMap(i, &p)

	// 	if err != nil {
	// 		fmt.Println("Got error from unmarshalling :")
	// 		fmt.Println(err.Error())
	// 		os.Exit(1)
	// 	}

	// 	if p.Age < maxAge+1 {
	// 		numItems++
	// 		fmt.Println("Name: ", p.Name)
	// 		fmt.Println("Age: ", p.Name)
	// 		fmt.Println("Breed: ", p.Breed)
	// 		fmt.Println()
	// 	}
	// }

	// fmt.Println("Found", numItems, "age less than or equal to ", maxAge, " with breed ", breed)

	puppyList := []*rpc.Puppy{}

	err2 := awsutils.UnmarshalListOfMaps(result.Items, &puppyList)

	if err2 != nil {
		fmt.Println(" Unmarshalling list of maps failed")
	}
	fmt.Printf("unmarshalled list %v \n", puppyList)

	return &rpc.PuppyList{
		Puppy: puppyList,
	}, nil
}
