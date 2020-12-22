package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bradfitz/slice"
	"math"

	"io/ioutil"
	"net/http"
)

type dataReq struct {
	Query string `json:"query"`
}

type Row struct {
	Month string
	Count int64
}

type RespData struct {
	Datarows []Row `json:"datarows"`
}

type esResp struct {
	Datarows [][]interface{} `json:"datarows"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	query := fmt.Sprintf("SELECT Month, COUNT(*) FROM border group by Month")
	fmt.Printf("query:%v", query)
	b, err := json.Marshal(&dataReq{Query: query})
	if err != nil {
		fmt.Printf("error in marshalling, err:%v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	resp, err := http.Post("https://search-border-es-tdsr3ykegaqmevpcmvhn57tizm.ap-south-1.es.amazonaws.com/_opendistro/_sql", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("error in requesting to es, err:%v", err)
		return events.APIGatewayProxyResponse{}, err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error in reading response, err:%v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Printf("result:%v", string(b))

	esResp := &esResp{}
	err = json.Unmarshal(b, esResp)
	if err != nil {
		fmt.Printf("error in unmarshall response, err:%v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	var rowsData []Row
	for _, dd := range esResp.Datarows {
		if len(dd) < 2 {
			continue
		}
		rowsData = append(rowsData, Row{
			Month: dd[0].(string)[:7],
			Count: int64(dd[1].(float64)),
		})
	}
	var maxNum float64
	var minNum float64
	maxNum = 1
	minNum = 1000000000
	for _, row := range rowsData {
		maxNum = math.Max(maxNum, float64(row.Count))
		minNum = math.Min(minNum, float64(row.Count))
	}
	fmt.Printf("min:%v, max:%v", minNum, maxNum)

	for i, row := range rowsData {
		count := row.Count
		rowsData[i].Count = int64(((float64(count) - minNum) * 100) / (maxNum - minNum))
		fmt.Printf("oldCount:%v, newCount:%v, month:%v", count, rowsData[i].Count, row.Month)
	}

	slice.Sort(rowsData[:], func(i, j int) bool {
		return rowsData[i].Month < rowsData[j].Month
	})
	b, err = json.Marshal(RespData{Datarows: rowsData})
	if err != nil {
		fmt.Printf("error in final marshall, err:%v", err)
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}
func main() {
	lambda.Start(handleRequest)
}
