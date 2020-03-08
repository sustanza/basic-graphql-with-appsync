package main

import (
	"basic-graphql-with-appsync/models"
	"basic-graphql-with-appsync/shared"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
	"math/rand"
	"os"
)

var CoffeeTableName = os.Getenv("COFFEE_TABLE_NAME")
var Roasts = []string{"Light", "Medium", "Dark"}
var Process = []string{"Machine Washed", "Dry Process", "Honey Process", "Wet Hulled", "Giling Basah", "Pulp Natural", "Wet Process"}
var Prefix = []string{"Organic", "Shade Grown", "High Mountain"}
var Postfix = []string{"Farm", "Coop", "Village"}

type RandomCoffeesPayload struct {
	Quantity int `json:"quantity"`
}

type RandomCoffeesResult struct {
	Coffees []models.Coffee `json:"coffees"`
}

func handle(payload RandomCoffeesPayload) (RandomCoffeesResult, error) {

	sess, err := shared.GetAWSSession()

	svc := dynamodb.New(sess)

	tableName := CoffeeTableName

	project := expression.NamesList(expression.Name("id"), expression.Name("name"), expression.Name("origin"), expression.Name("roast"))

	expr, err := expression.NewBuilder().WithProjection(project).Build()

	if err != nil {
		log.Printf("Got error building expression:")
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression: expr.Projection(),
		TableName:            aws.String(tableName),
	}

	result, err := svc.Scan(params)

	if err != nil {
		log.Printf("Query API call failed:%s", err.Error())
	}

	randomCoffeesResult := RandomCoffeesResult{}

	if len(result.Items) == 0 {
		log.Printf("Result returned by 0 items")
		return randomCoffeesResult, nil
	}

	var coffees []models.Coffee

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &coffees)

	if err != nil {
		log.Printf("Error unmarshaling results from dynamodb: %s", err.Error())
	}

	randomCoffees := inefficientKindaRandomResult(coffees, payload.Quantity)

	randomCoffeesResult.Coffees = randomCoffees

	log.Printf("Final output: %v", randomCoffeesResult)

	return randomCoffeesResult, nil
}

func inefficientKindaRandomResult(coffees []models.Coffee, quantity int) (randomCoffees []models.Coffee){

	for i := range coffees {
		j := rand.Intn(i + 1)
		coffees[i], coffees[j] = coffees[j], coffees[i]
	}

	randomCoffees = append([]models.Coffee{}, coffees[:quantity]...)

	return randomCoffees
}

func main() {
	lambda.Start(handle)
}
