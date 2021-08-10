package main

import (
	"fmt"
	"github.com/sokolovss/BNBsite/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan

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
		fmt.Println(err)
	}

	msg := mail.NewMSG()
	msg.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	msg.SetBody()

}
