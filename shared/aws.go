package shared

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)
func GetAWSSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{})

	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return sess, nil
}
