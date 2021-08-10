package main

import (
	"fmt"
	"github.com/sokolovss/BNBsite/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendEmail(msg)
		}
	}()

}

func sendEmail(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	msg := mail.NewMSG()
	msg.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		msg.SetBody(mail.TextHTML, m.Content)
	} else {
		fmt.Println("Using email template")
		data, err := ioutil.ReadFile(fmt.Sprintf("email-templates/%s", m.Template))
		if err != nil {
			errorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		msg.SetBody(mail.TextHTML, msgToSend)
	}

	err = msg.Send(client)
	if err != nil {
		errorLog.Println(err)
	} else {
		log.Printf("Email has been sent to %v", m.To)
	}

}
