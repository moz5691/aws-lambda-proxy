package awsutils

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

var DynamoTableName string

var DynamoClient = &atomicData{}

type atomicData struct {
	data   interface{}
	rwLock sync.RWMutex
}

func (m *atomicData) Get() interface{} {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	return m.data
}

func (m *atomicData) Set(data interface{}) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.data = data
}

func init() {

	// set dynamodb table name here
	DynamoTableName = "puppies"

	// region in AWS
	// dynamoConfig := &aws.Config{Region: aws.String("us-east-1")}

	// Check if dynamodb local is used or docker.
	// if (runtime.GOOS == "darwin") || (runtime.GOOS == "windows") {
	// 	dynamoConfig.Endpoint = aws.String("http://localhost:8000")
	// } else if os.Getenv("LISTEN_PORT") == "8080" {
	// 	dynamoConfig.Endpoint = aws.String("http://dynamo:8000")
	// }

	// sess, err := session.NewSession(dynamoConfig)

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	if (runtime.GOOS == "darwin") || (runtime.GOOS == "windows") {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		DynamoClient.Set(dynamodb.New(sess))
	} else {
		dynamoConfig := &aws.Config{Region: aws.String("us-east-1")}
		sess, err := session.NewSession(dynamoConfig)

		if err != nil {
			log.Println(err)
			return
		}
		DynamoClient.Set(dynamodb.New(sess))
	}

}

func GetDynamoClient() *dynamodb.DynamoDB {
	return DynamoClient.Get().(*dynamodb.DynamoDB)
}

// MarshalMap : aliasing of dynamodbattribute.MarshalMap
func MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		log.WithFields(log.Fields{
			"cause": err,
			"data":  fmt.Sprintf("%+v", in),
		}).Error("Error from marshaling data")
	}
	return av, err
}

// UnmarshalMap : aliasing of dynamodbattribute.UnmarshalMap
func UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	err := dynamodbattribute.UnmarshalMap(m, out)
	if err != nil {
		log.WithFields(log.Fields{
			"cause": err,
			"data":  fmt.Sprintf("%+v", m),
		}).Error("Error from unmarshalling")
	}
	return nil
}

// UnmarshalListOfMaps aliasing of dynmodbattribute.UnmarshalListOfMaps
func UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	err := dynamodbattribute.UnmarshalListOfMaps(l, out)
	if err != nil {
		log.WithFields(log.Fields{
			"cause": err,
			"data":  fmt.Sprintf("%+v", l),
		}).Error("Error from unmarshalling")
	}
	return nil
}

// GetItem wraps dynamodb.GetItem.
func GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	client := GetDynamoClient()
	input.TableName = &DynamoTableName
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.GetItemWithContext(ctx, input)
}

// DeleteItem wraps dynmodb.DeleteItem.
func DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	client := GetDynamoClient()
	input.TableName = &DynamoTableName
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.DeleteItemWithContext(ctx, input)
}

// UpdateItem wraps dynamodb.UpdateItem.
func UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	client := GetDynamoClient()
	input.TableName = &DynamoTableName
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.UpdateItemWithContext(ctx, input)
}

// PutItem wraps dynamodb.PutItem.
func PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	client := GetDynamoClient()
	input.TableName = &DynamoTableName
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.PutItemWithContext(ctx, input)
}

// Query wraps dynamodb.Query.
func Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	client := GetDynamoClient()
	input.TableName = &DynamoTableName
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.QueryWithContext(ctx, input)

}

func Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	client := GetDynamoClient()
	input.TableName = &DynamoTableName
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.ScanWithContext(ctx, input)
}
