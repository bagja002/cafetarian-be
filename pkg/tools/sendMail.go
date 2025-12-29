package tools

import (
	"fmt"
	"log"
	"project-kelas-santai/internal/config"

	"gopkg.in/gomail.v2"
)

const (
	SMTPHost   = "smtp.gmail.com"
	SMTPPort   = 587
	SenderName = "Cafe Santai <kelasantai.bootcamp@gmail.com>"
	AuthEmail  = "kelasantai.bootcamp@gmail.com"
)

// SendOrderSuccessEmail mengirim email notifikasi pesanan berhasil
func SendOrderSuccessEmail(toEmail, userName, orderID, message string) error {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	subject := fmt.Sprintf("Pesanan #%s Anda Sedang Diproses", orderID)
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f9f9f9; padding: 20px;">
			<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); overflow: hidden;">
				<div style="background-color: #4CAF50; padding: 20px; text-align: center; color: white;">
					<h2 style="margin: 0;">Terima Kasih!</h2>
				</div>
				<div style="padding: 30px;">
					<p style="font-size: 16px;">Halo <b>%s</b>,</p>
					<p style="font-size: 16px;">%s</p>
					<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
					<p style="font-size: 14px; color: #888;">Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi kami.</p>
				</div>
				<div style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 12px; color: #666;">
					<p>&copy; 2024 Cafe Santai. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, userName, message)

	m := gomail.NewMessage()
	m.SetHeader("From", SenderName)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	fmt.Println("Sending email to:", toEmail)
	fmt.Println("Email body:", body)
	fmt.Println("Email subject:", subject)
	fmt.Println("Email sender:", SenderName)
	fmt.Println("Email auth email:", AuthEmail)
	fmt.Println("Email app password:", cfg.Web.AppPassword)

	d := gomail.NewDialer(SMTPHost, SMTPPort, AuthEmail, cfg.Web.AppPassword)

	return d.DialAndSend(m)
}
