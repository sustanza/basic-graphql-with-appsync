package shared

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetAWSSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{})

	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return sess, nil
}
