package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Data struct {
	Month        string
	Cupcake      string
	Updated_time string
}

func main() {
	// read data from CSV file

	csvFile, err := os.Open("~/Downloads/multiTimeline.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	csvData, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var oneRecord Data
	var allRecords []Data

	for i, each := range csvData {
		if i == 0 {
			continue
		}
		oneRecord.Month = each[0]
		oneRecord.Cupcake = each[1]
		oneRecord.Updated_time = strconv.FormatInt(time.Now().Unix(), 10)
		allRecords = append(allRecords, oneRecord)
	}
	// converting to JSON
	jsondata, err := json.Marshal(allRecords)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(jsondata))

	//jsonFile, err := os.Create("data.json")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer jsonFile.Close()
	//jsonFile.Write(jsondata)

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	for _, data := range allRecords {
		av, err := dynamodbattribute.MarshalMap(data)
		if err != nil {
			fmt.Println("Got error marshalling map:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// Create item in table Movies
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Data"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func GenRand() string {
	rand.Seed(time.Now().UnixNano())
	n := 2020 - 2004 + 1
	num := rand.Intn(n)
	year := 2004 + num

	month := rand.Intn(12) + 1

	monthString := fmt.Sprintf("%d-%02d", year, month)
	return monthString
}


