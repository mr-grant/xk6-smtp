import smtp from 'k6/x/smtp';


export default function () {
    smtp.sendMail(
        "smtp.gmail.com",
        "587",
        "sender@gmail.com",
        "senderPassword",
        ["recipient1@gmail.com", "recipient2@gmail.com"],
        "Test subject",
        "Test message"
    )
}
