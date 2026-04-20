package Utils

import (
	"fmt"
	"net/smtp"
	"os" 
	"net"  
	"crypto/tls"
)

func SendPasswordResetEmail(userEmail, resetToken string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_APP_PASSWORD") 
	fmt.Println("EMAIL:", from)
	fmt.Println("PASS LENGTH:", len(pass))
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 🌐 Change this to your actual frontend URL (e.g., your React App URL)
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password/%s", resetToken)

	// 📧 Define Headers and Body
	header := make(map[string]string)
	header["From"] = from
	header["To"] = userEmail
	header["Subject"] = "Password Reset Request"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	// Construct the message headers
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// 🖌️ HTML Body (Professional Look)
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px;">
				<h2 style="color: #007bff;">Reset Your Password</h2>
				<p>Hello,</p>
				<p>We received a request to reset the password for your account. Click the button below to choose a new password:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; font-weight: bold; display: inline-block;">
						Reset Password
					</a>
				</div>
				<p>This link will expire in 15 minutes. If you did not request this change, please ignore this email.</p>
				<p>Thanks,<br>The Nivas Team</p>
			</div>
		</body>
		</html>
	`, resetLink)

	fullMessage := message + "\r\n" + body

	// 🚀 Send the email
	auth := smtp.PlainAuth("", from, pass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{userEmail}, []byte(fullMessage))
} 

func SendWelcomeEmail(userEmail, tempPassword string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_APP_PASSWORD")

	fmt.Println("EMAIL:", from)
	fmt.Println("PASS LENGTH:", len(pass))
	
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 📧 Headers
	header := make(map[string]string)
	header["From"] = from
	header["To"] = userEmail
	header["Subject"] = "Welcome"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Welcome!</h2>
		<p>Username: %s</p>
		<p>Password: %s</p>
	</body>
	</html>
	`, userEmail, tempPassword)

	fullMessage := message + "\r\n" + body

	// 🔥 CONNECT manually
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}

	// 🔐 STARTTLS (IMPORTANT)
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	// 🔑 AUTH
	auth := smtp.PlainAuth("", from, pass, smtpHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// 📤 SEND
	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(userEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fullMessage))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}