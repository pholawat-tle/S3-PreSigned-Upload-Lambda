package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	lambda.Start(HandleEvent)
}

func HandleEvent(e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set up Lambda Function
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Create S3 Client
	client := s3.NewPresignClient(s3.NewFromConfig(cfg))
	bucketName := "your-s3-bucket-name"
	// Get file name from Request
	var RequestBody interface{}
	json.Unmarshal([]byte(e.Body), &RequestBody)
	key := RequestBody.(map[string]interface{})["Key"].(string)
	// Generate Presigned URL for upload
	params := s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}

	PresignedRequest, err := client.PresignPutObject(context.TODO(), &params)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"PresignedURL": %s`, PresignedRequest.URL),
	}, nil
}
