package Utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendPasswordResetEmail sends an email with a link to reset the user's password
func SendPasswordResetEmail(userEmail, resetToken string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_APP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	
	// Default to 587 if PORT is not set or invalid
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if smtpPort == 0 {
		smtpPort = 587
	}

	resetLink := fmt.Sprintf("http://localhost:3000/reset-password/%s", resetToken)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Password Reset Request")

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

	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, from, pass)

	return d.DialAndSend(m)
}

// SendWelcomeEmail sends the temporary login credentials to new users
func SendWelcomeEmail(userEmail, tempPassword string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_APP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	
	
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	
	
	if smtpPort == 0 {
		smtpPort = 587
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Welcome")

	body := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif; line-height: 1.6;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #eee; border-radius: 8px;">
			<h2 style="color: #28a745;">Welcome to the Team!</h2>
			<p>Your account has been created. Use the temporary credentials below to log in:</p>
			<div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
				<p><strong>Username/Email:</strong> %s</p>
				<p><strong>Temporary Password:</strong> <code style="color: #d63384;">%s</code></p>
			</div>
			<p>For security, you will be required to change this password upon your first login.</p>
			
		</div>
	</body>
	</html>
	`, userEmail, tempPassword)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, from, pass)

	return d.DialAndSend(m)
}