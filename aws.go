package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ses"
)

func SendEmail(notification Notification) {

	svc := ses.New(&aws.Config{Region: "us-east-1"})

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			// BCCAddresses: []*string{
			// 	aws.String(""),
			// },
			ToAddresses: []*string{
				aws.String(notification.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				HTML: &ses.Content{
					Data:    aws.String(notification.Message),
					Charset: aws.String("utf-8"),
				},
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String(notification.Subject),
				Charset: aws.String("utf-8"),
			},
		},
		Source: aws.String("alert@chrisbenson.com"),
		ReplyToAddresses: []*string{
			aws.String("alert@chrisbenson.com"),
		},
		ReturnPath: aws.String("alert@chrisbenson.com"),
	}
	resp, err := svc.SendEmail(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS Error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, The SDK should alwsy return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(resp))
}
