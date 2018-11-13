package main

import (
    "fmt"
    "os"
    "time"
    "log"
    "net/http"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
    Name string`json:"Name"`
    Message string`json:"Message"`
}

var dynamodbClient *dynamodb.DynamoDB

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
    dynamodbClient = dynamodb.New(sess)

    http.HandleFunc("/write", writeHandler)
    http.HandleFunc("/read", readHandler)
    log.Fatal(http.ListenAndServe(":80", nil))
}


// Read the SERVICE index from the TABLE_NAME
func readHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    
    result, err := dynamodbClient.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
        Key: map[string]*dynamodb.AttributeValue{
            "Name": {
                S: aws.String("SERVICE"),
            },
        },
    })
    
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    
    item := Item{}
    
    err = dynamodbattribute.UnmarshalMap(result.Item, &item)

    // Return the Message
    fmt.Fprintf(w, item.Message)
}

// Write to the SERVICE index 
func writeHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)

    message := r.URL.Query().Get("message")
    if message == "" {
        message = time.Now().Format(time.RFC850)
    }
    item := Item{
        // Name: os.Getenv("APP_NAME"),
        Name: "SERVICE",
        Message: message,
    }
    av, err := dynamodbattribute.MarshalMap(item)

    input := &dynamodb.PutItemInput{
        Item: av,
        TableName: aws.String(os.Getenv("TABLE_NAME")),
    }
    
    _, err = dynamodbClient.PutItem(input)

    if err != nil {
        fmt.Println("Got error calling PutItem:")
        fmt.Println(err.Error())
        os.Exit(1)
    }
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}