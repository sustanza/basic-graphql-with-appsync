package main

import (
	"basic-graphql-with-appsync/models"
	"log"
	"math/rand"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// defines the payload expected from the resolver request
type RandomCoffeesPayload struct {
	Quantity int `json:"quantity"`
}

// defines the result the resolver expects as a response
type RandomCoffeesResult struct {
	Coffees []models.Coffee `json:"coffees"`
}

// Handle: handles the lambda request for random coffee resources
func handle(payload RandomCoffeesPayload) (RandomCoffeesResult, error) {

	// defines the result for future use
	randomCoffeesResult := RandomCoffeesResult{}

	// creates a new AWS session
	sess, err := session.NewSession(&aws.Config{})

	// log an error if the session cannot be created and return back an empty result
	if err != nil {
		log.Printf("%v", err)
		return randomCoffeesResult, err
	}

	// creates a new dynamodb service
	svc := dynamodb.New(sess)

	// grabs the coffee table name from the env
	tableName := os.Getenv("COFFEE_TABLE_NAME")

	// creates a namelist to be used as part of a expression builder defining the values to return back from the db
	project := expression.NamesList(expression.Name("id"), expression.Name("name"), expression.Name("origin"), expression.Name("roast"))

	// creates a new expression builder
	expr, err := expression.NewBuilder().WithProjection(project).Build()

	// if there is an issue building the expression we log an error and return back an empty result
	if err != nil {
		log.Printf("Got error building expression:")
		return randomCoffeesResult, err
	}

	// defines all the parameters for scanning dynamodb
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// performs a scan on dynamodb to grab coffee resources
	result, err := svc.Scan(params)

	// if there is an issue with the scan we log it and return back an empty result
	if err != nil {
		log.Printf("Query API call failed:%s", err.Error())
		return randomCoffeesResult, nil
	}

	// if there is no coffee resources in the result we log it and return back an empty result
	if len(result.Items) == 0 {
		log.Printf("Result returned by 0 items")
		return randomCoffeesResult, nil
	}

	// define a new slice to store our coffee resources
	var coffees []models.Coffee

	// unmarshal the list of resources so we can go use them
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &coffees)

	// if there is an issue with the unmarshal we log it and return back an empty result
	if err != nil {
		log.Printf("Error unmarshaling results from dynamodb: %s", err.Error())
		return randomCoffeesResult, nil
	}

	// gets some random coffee resources very inefficiently
	randomCoffees := inefficientKindaRandomResult(coffees, payload.Quantity)

	// adds the coffee resources to our result
	randomCoffeesResult.Coffees = randomCoffees

	// prints out the final result for logging
	log.Printf("Final output: %v", randomCoffeesResult)

	// returns the result
	return randomCoffeesResult, nil
}

// inefficientKindaRandomResult: inefficiently returns back kinda random results of coffee resources
func inefficientKindaRandomResult(coffees []models.Coffee, quantity int) (randomCoffees []models.Coffee) {

	// check if the quantity is defined, if it doesn't it just set the quantity to 1
	if quantity == 0 {
		quantity = 1
	}

	// iterates over a slice of coffees
	for i := range coffees {
		// generates a random int to be used in scrambling some coffee orders
		j := rand.Intn(i + 1)
		coffees[i], coffees[j] = coffees[j], coffees[i]
	}

	// creates new coffee slices out of the previously scrambled slices
	randomCoffees = append([]models.Coffee{}, coffees[:quantity]...)

	// returns back the inefficiently randomized kinda random results haha
	return randomCoffees
}

func main() {
	lambda.Start(handle)
}
