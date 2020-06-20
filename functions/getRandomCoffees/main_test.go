package main

import (
	"basic-graphql-with-appsync/models"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/brianvoe/gofakeit/v4"
)

// these are a set of string slices to generate
var Roasts = []string{"Light", "Medium", "Dark"}
var Process = []string{"Machine Washed", "Dry Process", "Honey Process", "Wet Hulled", "Giling Basah", "Pulp Natural", "Wet Process"}
var Prefix = []string{"Organic", "Shade Grown", "High Mountain"}
var Postfix = []string{"Farm", "Coop", "Village"}

func TestGetRandomCoffeesResolver(t *testing.T) {
	randomCoffeesPayload := RandomCoffeesPayload{
		Quantity: 6,
	}
	data, err := handle(randomCoffeesPayload)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if data.Coffees == nil {
		t.Errorf("Must return with coffees")
	}
}

func TestSeedData(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{})

	if err != nil {
		t.Errorf("%s", err.Error())
		os.Exit(1)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := os.Getenv("COFFEE_TABLE_NAME")

	coffees := [100]*models.Coffee{}

	gofakeit.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {

		coffees[i] = &models.Coffee{
			Id: gofakeit.UUID(),
			Name: fmt.Sprintf("%s %s %s %s",
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
