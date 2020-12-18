package game

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"math/rand"
	"time"
)

type DDBMutex struct {
	TimeStamp int64
	NodeID string
}

type Score struct {
	Name string
	Score int64
}

var (
	mutexTableName = "scribbleScoreMutex"
)

func lockDDBMutex(nodeId string) error{
	fmt.Println("locking")
	timeStamp := time.Now().Unix()
	counter := 1
	for time.Now().Unix() - timeStamp < 120 {
		fmt.Println(counter)
		counter++
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-west-2")},
		)
		if err != nil {
			return err
		}

		svc := dynamodb.New(sess)
		result, err := svc.Scan(&dynamodb.ScanInput{
			TableName: aws.String(mutexTableName),
		})
		if err != nil {
			return err
		}

		items := result.Items
		fmt.Println("mutex value items")
		fmt.Println(items)
		fmt.Println("end")
		existingMutex := DDBMutex{}
		if items != nil && len(items) > 0 {
			err = dynamodbattribute.UnmarshalMap(items[0], &existingMutex)
			fmt.Println("err")
			fmt.Println(err)
		}
		// no items or mutex expired after 5 minutes
		if items == nil || time.Now().Unix() - existingMutex.TimeStamp > 300 {
			av, err := dynamodbattribute.MarshalMap(&DDBMutex{
				TimeStamp: time.Now().Unix(),
				NodeID:    nodeId,
			})
			if err != nil {
				return err
			}

			input := &dynamodb.PutItemInput{
				Item: av,
				TableName: aws.String(mutexTableName),
			}

			if _, err = svc.PutItem(input); err != nil {
				return err
			}
		}

		if existingMutex.NodeID == nodeId {
			return nil
		}
		fmt.Println("existing mutex node id: " + existingMutex.NodeID)
		fmt.Println("current nodeID: " + nodeId)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
	return errors.New("Timeout getting ddb mutex")
}

func UnlockDDBMutex(nodeId string) error {
	fmt.Println("unlocking")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return err
	}
	svc := dynamodb.New(sess)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"NodeID": {
				S: aws.String(nodeId),
			},
		},
		TableName: aws.String(mutexTableName),
	}

	_, err = svc.DeleteItem(input)
	return err
}