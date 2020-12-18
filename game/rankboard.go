package game

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strings"
	"time"
)

//TODO
// Implement rank over score
// Rank is a map of player name to the scoreMap score
// Rank is calculated by how many people a play bests in all previous games:
// Winner of a game of 5 will get a scoreMap score of 4 where last place will get a scoreMap score of 0
var (
	scoreBoard []scoreEntry
	scoreTableName = "scribbleScore"
)

func InitializeRank() {
	scoreBoard = []scoreEntry{}
	scoreBoard = append(scoreBoard, scoreEntry{Name: "snowie", Score: 99})
	scoreBoard = append(scoreBoard, scoreEntry{Name: "xuechenf", Score: 100})
	go keepSyncScore()
}

type scoreEntry struct {
	Name string
	Score int
}

func PrintScore() string{
	builder := strings.Builder{}
	for k, v := range scoreBoard {
		builder.WriteString(fmt.Sprintf("%s -> %d <br />\n\r", k, v)) }
	return builder.String()
}

func keepSyncScore() error {
	for {
		fmt.Println("syncing")
		if err := syncScoreBoard(); err != nil {

			fmt.Println("sync board error:")
			fmt.Println(err)
			fmt.Println("end")
			return err
		}
		time.Sleep(time.Duration(60) * time.Second)
	}
}

func syncScoreBoard() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(scoreTableName),
	}
	result, err := svc.Scan(params)
	fmt.Print("scan result")
	fmt.Print(result)

	if err != nil {
		return err
	}
	if result.Items != nil {
		return dynamodbattribute.UnmarshalListOfMaps(result.Items, &scoreBoard)
	}
	return nil
}

func UpdatePlayerScores(players []*Player) error{
	fmt.Println("lock result:")
	fmt.Println(lockDDBMutex())
	fmt.Println("end")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	fmt.Println(len(players))
	fmt.Println(len(players))
	fmt.Println(len(players))

	for _, player := range players {
		fmt.Println("Update player score")
		name := player.Name
		result, err := svc.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(scoreTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Name" : {
					S: aws.String(name),
				},
			},
		})
		if err != nil {
			fmt.Println("error on getItem")
			return err
		}
		fmt.Println("passed getitem")
		item := scoreEntry{}
		if result.Item != nil {
			err = dynamodbattribute.UnmarshalMap(result.Item, &item)
		}
		fmt.Println(name)
		item.Name = name
		fmt.Println(item.Score)
		item.Score += player.Score
		fmt.Println(item.Score)

		// Update the table
		av, err := dynamodbattribute.MarshalMap(&item)
		if err != nil {
			return err
		}
		input := &dynamodb.PutItemInput{
			Item: av,
			TableName: aws.String(scoreTableName),
		}
		if _, err = svc.PutItem(input); err != nil {
			return err
		}
	}
	return UnlockDDBMutex()
}
