package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Notification struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type SNSMessage struct {
	Message string `json:"Message"`
}

func sendEmail(
	to string,
	subject string,
	body string,
) error {

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPassword,
		smtpHost,
	)

	message := []byte(
		fmt.Sprintf(
			"Subject: %s\r\n\r\n%s",
			subject,
			body,
		),
	)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpUser,
		[]string{to},
		message,
	)
}

func Handler(
	ctx context.Context,
	sqsEvent events.SQSEvent,
) error {

	for _, record := range sqsEvent.Records {

		var snsMsg SNSMessage

		if err := json.Unmarshal(
			[]byte(record.Body),
			&snsMsg,
		); err != nil {

			log.Println(err)
			continue
		}

		var notification Notification

		if err := json.Unmarshal(
			[]byte(snsMsg.Message),
			&notification,
		); err != nil {

			log.Println(err)
			continue
		}

		err := sendEmail(
			notification.Email,
			notification.Subject,
			notification.Message,
		)

		if err != nil {
			log.Println("Error enviando correo:", err)
			continue
		}

		log.Printf(
			"Correo enviado a %s",
			notification.Email,
		)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}