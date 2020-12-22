package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type DataEs struct {
	Month        string
	Cupcake      string
	Updated_time string
	EventType    string
}

func handler(e events.DynamoDBEvent) error {
	var item map[string]events.DynamoDBAttributeValue
	fmt.Println("Beginning ES Sync")

	for _, v := range e.Records {
		switch v.EventName {
		case "INSERT":
			fallthrough
		case "MODIFY":
			tableName := strings.Split(v.EventSourceArn, "/")[1]
			fmt.Printf("tableName:%v", tableName)
			item = v.Change.NewImage
			data := DataEs{
				Month:        item["Month"].String(),
				Cupcake:      item["Cupcake"].String(),
				Updated_time: item["Updated_time"].String(),
				EventType:    v.EventName,
			}

			fmt.Printf("item:%+v", item)

			b, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("error in marshalling, err:%v", err)
				continue
			}
			fmt.Printf("req:%v", string(b))

			resp, err := http.Post("https://search-border-es-tdsr3ykegaqmevpcmvhn57tizm.ap-south-1.es.amazonaws.com/border/lambda", "application/json", bytes.NewBuffer(b))
			if err != nil {
				fmt.Printf("error in inserting in es, err:%v", err)
				continue
			}
			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("error in reading response, err:%v", err)
				continue
			}
			fmt.Printf("response:%v", string(b))
		default:
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
