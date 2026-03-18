package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/smtp", new(SMTP))
}

type SMTP struct{}

func (*SMTP) SendMail(host string, port string, sender string, senderPassword string, recipient string, title string, message string) {
	// Configuration
	smtpHost := host
	smtpPort := port
	username := sender
	password := senderPassword
	from := sender
	to := []string{recipient}

	// 1. Connection (plain, without TLS)
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer conn.Close()

	// 2. Create SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Printf("Client creation error: %v\n", err)
		return
	}
	defer client.Close()

	// 3. Upgrade to TLS via STARTTLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // For testing only!
		ServerName:         smtpHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Printf("STARTTLS error: %v\n", err)
		return
	}

	// 4. Authentication
	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		fmt.Printf("Authentication error: %v\n", err)
		return
	}

	// 5. Send email
	if err = client.Mail(from); err != nil {
		fmt.Printf("Mail error: %v\n", err)
		return
	}

	for _, rcpt := range to {
		if err = client.Rcpt(rcpt); err != nil {
			fmt.Printf("Rcpt error: %v\n", err)
			return
		}
	}

	writer, err := client.Data()
	if err != nil {
		fmt.Printf("Data error: %v\n", err)
		return
	}

	// Build email
	emailMsg := "To: " + to[0] + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + title + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + message +
		"\r\n"

	_, err = writer.Write([]byte(emailMsg))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("Close error: %v\n", err)
		return
	}
}
