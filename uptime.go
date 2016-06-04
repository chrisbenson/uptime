package main

import (
	// "encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

var idChan = make(chan uuid.UUID)
var statusChan = make(chan Status)
var Websites []Website
var lastStatus Status
var statusBeforeLast Status

type Website struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Interval int    `json:"interval"` // in seconds
	Email    string `json:"email"`
}

type Status struct {
	Id        uuid.UUID `json:"id"`
	WebsiteId int       `json:"website-id"`
	Time      time.Time `json:"time"`
	Url       string    `json:"url"`
	Code      int       `json:"code,omitempty"`
	Error     error     `json:"error,omitempty"`
}

type Notification struct {
	Id      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Subject string    `json:"subject"`
	Message string    `json:"message"`
	Url     string    `json:"url"`
}

func UuidProvider() {
	for {
		idChan <- uuid.NewV1()
	}

}

func MonitorWebsite(website Website) {

	var lastStatus Status
	var statusBeforeLast Status

	for {
		code, err := GetStatusCode(website.Url)
		status := CreateStatus(website.Id, website.Url, code, err)

		if status.Code >= 400 &&
			status.Code == lastStatus.Code &&
			status.Code != statusBeforeLast.Code {
			CreateFailureNotification(website)
		} else if status.Code >= 400 &&
			status.Code == lastStatus.Code &&
			&statusBeforeLast == nil {
			fmt.Println("Sending failure notification about:\n" + website.Url)
			CreateFailureNotification(website)
		}

		statusBeforeLast = lastStatus
		lastStatus = status

		statusChan <- status
		time.Sleep(time.Duration(website.Interval) * time.Second)
	}
}

func GetStatusCode(url string) (int, error) {
	res, err := http.Get(url)
	return res.StatusCode, err
}

func CreateStatus(websiteId int, url string, code int, err error) Status {
	return Status{
		Id:        <-idChan,
		WebsiteId: websiteId,
		Time:      time.Now().UTC(),
		Url:       url,
		Code:      code,
		Error:     err,
	}
}

func CreateNotification(email string, subject string, message string, url string) Notification {
	return Notification{
		Id:      <-idChan,
		Email:   email,
		Subject: subject,
		Message: message,
		Url:     url,
	}
}

func SendNotification(notification Notification) {
	SendEmail(notification)
}

func CreateFailureNotification(website Website) {
	subject := "Test email from Chris Benson. Please ignore."
	message := "The website at this URL is currently not available: \n" + website.Url
	notification := CreateNotification(website.Email, subject, message, website.Url)
	SendNotification(notification)

}
