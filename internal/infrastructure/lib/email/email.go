package email

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/email"
)

type SMTPSender struct {
	email.EmailSender
	log *slog.Logger
	cfg *config.Config
}

func NewSMTPSender(log *slog.Logger, cfg *config.Config) *SMTPSender {
	return &SMTPSender{
		log: log,
		cfg: cfg,
	}
}

func (s *SMTPSender) SendEmail(to string, subject string, body string) error {
	subject = fmt.Sprintf("Subject: %s\n", subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"

	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				/* Общие стили */
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					color: #333333;
					margin: 0;
					padding: 0;
				}
				.container {
					width: 80%%;
					margin: 20px auto;
					padding: 20px;
					background-color: #ffffff;
					border-radius: 10px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.header {
					text-align: center;
					padding-bottom: 20px;
					border-bottom: 1px solid #dddddd;
				}
				.header h1 {
					color: #4CAF50;
					font-size: 28px;
					margin: 0;
				}
				.content {
					padding: 20px;
				}
				.content p {
					font-size: 16px;
					line-height: 1.6;
				}
				.button {
					display: inline-block;
					background-color: #4CAF50;
					color: #ffffff;
					padding: 10px 20px;
					text-decoration: none;
					border-radius: 5px;
					font-size: 16px;
				}
				.footer {
					text-align: center;
					margin-top: 20px;
					font-size: 12px;
					color: #777777;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Смена пароля</h1>
				</div>
				<div class="content">
					<p>Привет,</p>
					<p>Твой временный код - %s</p>
					<p>Если ты не запрашивал восстановление пароля - срочно выйди из всех аккаунтов и смени пароль!</p>
					<!-- <a href="#" class="button">Красивая кнопка</a> -->
				</div>
				<div class="footer">
					<p>Спасибо за внимание!<br>Вы получили это письмо в целях тестирования.</p>
				</div>
			</div>
		</body>
		</html>`, body)

	msg := []byte(subject + mime + "\n" + htmlBody)

	addr := s.cfg.SMTPGmailHost + ":" + s.cfg.SMTPGmailPort

	auth := smtp.PlainAuth("", s.cfg.GmailEmail, s.cfg.SMTPPassword, s.cfg.SMTPGmailHost)
	fmt.Println(addr)
	fmt.Println(s.cfg.GmailEmail, s.cfg.SMTPPassword, s.cfg.SMTPGmailHost)
	err := smtp.SendMail(addr, auth, s.cfg.GmailEmail, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("ошибка при отправке письма: %v", err)
	}

	return nil
}
func (s *SMTPSender) SendTestEmail(to string, subject string, body string) error {
	subject = fmt.Sprintf("Subject: %s\n", subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"

	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				/* Общие стили */
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					color: #333333;
					margin: 0;
					padding: 0;
				}
				.container {
					width: 80%%;
					margin: 20px auto;
					padding: 20px;
					background-color: #ffffff;
					border-radius: 10px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.header {
					text-align: center;
					padding-bottom: 20px;
					border-bottom: 1px solid #dddddd;
				}
				.header h1 {
					color: #4CAF50;
					font-size: 28px;
					margin: 0;
				}
				.content {
					padding: 20px;
				}
				.content p {
					font-size: 16px;
					line-height: 1.6;
				}
				.button {
					display: inline-block;
					background-color: #4CAF50;
					color: #ffffff;
					padding: 10px 20px;
					text-decoration: none;
					border-radius: 5px;
					font-size: 16px;
				}
				.footer {
					text-align: center;
					margin-top: 20px;
					font-size: 12px;
					color: #777777;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Смена пароля</h1>
				</div>
				<div class="content">
					<p>Привет,</p>
					<p>Твой временный код - %s</p>
					<p>Если ты не запрашивал восстановление пароля - срочно выйди из всех аккаунтов и смени пароль!</p>
					<!-- <a href="#" class="button">Красивая кнопка</a> -->
				</div>
				<div class="footer">
					<p>Спасибо за внимание!<br>Вы получили это письмо в целях тестирования.</p>
				</div>
			</div>
		</body>
		</html>`, body)

	msg := []byte(subject + mime + "\n" + htmlBody)
	err := smtp.SendMail(s.cfg.SMTPHost+":"+s.cfg.SMTPPort, nil, s.cfg.OurEmail, []string{to}, msg)
	if err != nil {
		s.log.Error("Ошибка при отправке письма:", err)
		return err
	}
	s.log.Info("Письмо успешно отправлено!")
	return nil
}
