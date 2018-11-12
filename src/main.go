package main

import (
    "fmt"
    "os"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
    Name string`json:"Name"`
    LastWrite string`json:"LastWrite"`
}

func main() {

    if len(os.Getenv("APP_NAME")) == 0 {
        fmt.Println("APP_NAME environment var is missing")
        os.Exit(1)
    }
    if len(os.Getenv("AWS_REGION")) == 0 {
        fmt.Println("AWS_REGION environment var is missing")
        os.Exit(1)
    }
    if len(os.Getenv("AWS_ACCESS_KEY_ID")) == 0 {
        fmt.Println("AWS_ACCESS_KEY_ID environment var is missing")
        os.Exit(1)
    }
    if len(os.Getenv("AWS_SECRET_ACCESS_KEY")) == 0 {
        fmt.Println("AWS_SECRET_ACCESS_KEY environment var is missing")
        os.Exit(1)
    }
    if len(os.Getenv("TABLE_NAME")) == 0 {
        fmt.Println("TABLE_NAME environment var is missing")
        os.Exit(1)
    }
    
    // Initialize a session
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION"))},
    )

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    for {
        item := Item{
            Name: os.Getenv("APP_NAME"),
            LastWrite: time.Now().Format(time.RFC850),
        }
        av, err := dynamodbattribute.MarshalMap(item)
    
        input := &dynamodb.PutItemInput{
            Item: av,
            TableName: aws.String(os.Getenv("TABLE_NAME")),
        }
        
        _, err = svc.PutItem(input)
    
        if err != nil {
            fmt.Println("Got error calling PutItem:")
            fmt.Println(err.Error())
            os.Exit(1)
        }

        time.Sleep(5 * time.Second)
    }
}
