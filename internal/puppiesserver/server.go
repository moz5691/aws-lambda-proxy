package puppiesserver

import (
	"context"
	"fmt"

	// "github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	// "math/rand"

	"github.com/moz5691/lambda-proxy/internal/awsutils"
	rpc "github.com/moz5691/lambda-proxy/rpc/puppies"
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

func (s *Server) GetById(ctx context.Context, req *rpc.GetByIdReq) (*rpc.Puppy, error) {
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
