package smtp

import (
	"fmt"
	"net/smtp"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/smtp", new(SMTP))
}

type SMTP struct{}

func (*SMTP) SendMail(host string, port string, sender string, password string, recipient string, title string, message string) {
	// 📋 Конфигурация
	smtpHost := host
	smtpPort := port
	username := sender
	password := password
	from := recipient
	to := []string{recipient}

	// 1️⃣ Подключение (обычное, без TLS)
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		fmt.Printf("❌ Ошибка подключения: %v\n", err)
		return
	}
	defer conn.Close()

	// 2️⃣ Создание SMTP клиента
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Printf("❌ Ошибка создания клиента: %v\n", err)
		return
	}
	defer client.Close()

	// 3️⃣ Апгрейд до TLS через STARTTLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // ⚠️ Только для тестов!
		ServerName:         smtpHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Printf("❌ Ошибка STARTTLS: %v\n", err)
		return
	}

	// 4️⃣ Аутентификация
	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		fmt.Printf("❌ Ошибка аутентификации: %v\n", err)
		return
	}

	// 5️⃣ Отправка письма
	if err = client.Mail(from); err != nil {
		fmt.Printf("❌ Ошибка Mail: %v\n", err)
		return
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			fmt.Printf("❌ Ошибка Rcpt: %v\n", err)
			return
		}
	}

	writer, err := client.Data()
	if err != nil {
		fmt.Printf("❌ Ошибка Data: %v\n", err)
		return
	}

	// Формирование письма
	message := []byte("To: " + to[0] + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + title + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + message +
		"\r\n")

	_, err = writer.Write(message)
	if err != nil {
		fmt.Printf("❌ Ошибка записи: %v\n", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("❌ Ошибка закрытия: %v\n", err)
		return
	}
}
