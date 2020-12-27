package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GenRand() string {
	rand.Seed(time.Now().UnixNano())
	n := 2020 - 2004 + 1
	num := rand.Intn(n)
	year := 2004 + num

	month := rand.Intn(12) + 1

	monthString := fmt.Sprintf("%d-%02d", year, month)
	return monthString
}

func handler(_ interface{}) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
		for i := 0; i < 100; i++ {
			tableName := "Data"
			monthRand := GenRand()
			updateTime := strconv.FormatInt(time.Now().Unix(), 10)

			input := &dynamodb.UpdateItemInput{
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":Updated_time": {
						S: aws.String(updateTime),
					},
				},
				TableName: aws.String(tableName),
				Key: map[string]*dynamodb.AttributeValue{
					"Month": {
						S: aws.String(monthRand),
					},
				},
				ReturnValues:     aws.String("UPDATED_NEW"),
				UpdateExpression: aws.String("set Updated_time = :Updated_time"),
			}

			_, err := svc.UpdateItem(input)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			fmt.Println(monthRand)
			fmt.Println(updateTime)
			time.Sleep(time.Second)
		}
		return nil

}

func main () {
	lambda.Start(handler)
}
