package main

import (
	"basic-graphql-with-appsync/models"
	"basic-graphql-with-appsync/shared"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/brianvoe/gofakeit/v4"
	"os"
	"testing"
	"time"
)

func TestGetRandomCoffeesResolver(t *testing.T) {
	randomCoffeesPayload := RandomCoffeesPayload{
		Quantity: 3,
	}
	data, err := handle(randomCoffeesPayload)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if data.coffees == nil {
		t.Errorf("Must return with coffees")
	}
}

func TestSeedData(t *testing.T) {
	sess, err := shared.GetAWSSession()

	if err != nil {
		t.Errorf("%s", err.Error())
		os.Exit(1)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := CoffeeTableName

	coffees := [100]*models.Coffee{}

	gofakeit.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {

		coffees[i] = &models.Coffee{
			Id:     gofakeit.UUID(),
			Name:   fmt.Sprintf("%s %s %s %s",
				Prefix[gofakeit.Number(0, len(Prefix)-1)],
				Process[gofakeit.Number(0, len(Process)-1)],
				gofakeit.LastName(),
				Postfix[gofakeit.Number(0, len(Postfix)-1)],
			),
			Origin: gofakeit.Country(),
			Roast:  Roasts[gofakeit.Number(0, len(Roasts)-1)],
		}
	}

	for i := 0; i < 100; i++ {
		av, err := dynamodbattribute.MarshalMap(coffees[i])
		if err != nil {
			t.Errorf("%s", err.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			t.Errorf("%s", err.Error())
			os.Exit(1)
		}
	}
}
