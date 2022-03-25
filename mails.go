package main

import (
	Config "click-al-vet/config"
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

var htmlContent = "<div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>!Hola ¡Javier gil!, los detalles de tu cita son los siguientes: </span><br/><br/><table style='border-collapse: collapse; width:100%; border: 1px solid black;'  ><tbody><tr><td style='border: 1px solid black' ><b>Diagnostico:</b></td><td style='border: 1px solid black'  >Infeccion debida a coronavirus, sin otra especificacion</td></tr><tr><td style='border: 1px solid black'  ><b>Observaciones:</b></td><td style='border: 1px solid black'  >Se hicieron diferentes pruebas y se determino que el diagnostico es debido a ...</td></tr></tbody></table><br/><span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Los medicamentos recetados son los siguientes:</span><br/><br/><table style='border-collapse: collapse; width:100%; border: 1px solid black;'  ><thead><tr><th style='color: #54ace2;font-size: 16px;font-weight: bold;'>Medicamento</th><th  style='color: #54ace2;font-size: 16px;font-weight: bold;'>Presentación</th><th  style='color: #54ace2;font-size: 16px;font-weight: bold;'>Posología</th><th  style='color: #54ace2;font-size: 16px;font-weight: bold;'>Duración</th></tr></thead><tbody><tr><td>dsd</td><td>dsd</td><td>dsd</td><td>dsd</td></tr></tbody></table></div></div>"

func testEmail() {

	var config = Config.Config{}
	config.Read()

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	//m.SetHeader("To", "solucionesitecnologia@gmail.com")
	m.SetHeader("To", "ventas.javc@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Bienvenido a clickal medic, confirma tu contraseña")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContent)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendResetPasswordEmail(token string, mail string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = " <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Habilita tu usuario con este </span><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/recover-password?tokenizer=" + token + "' >Enlace</a></div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Saludos de clickal medic, confirma tu contraseña")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendConfirmationEmail(token string, mail string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = " <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Habilita tu usuario con este </span><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/confirm-account?tokenizer=" + token + "' >Enlace</a></div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Bienvenido a clickal medic, confirma tu usuario")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendForgotPasswordEmail(token string, mail string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = " <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Habilita tu usuario con este </span><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/recover-password?tokenizer=" + token + "' >Enlace</a></div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Recupera tu usario con el siguiente enlace")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendAppointmentConfirmationEmail(token string, mail string, appointment string, doctorName string, appointmentDate string, appointmentHour string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = "  <div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div><span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>El doctor </span><span style='color: #54ace2;font-size: 20px;font-weight: bold;'> " + doctorName + " </span><span style='color: #0f76b0;font-size: 20px;font-weight: bold;'> agendo su cita para:</span><span style='color: #54ace2;font-size: 20px;font-weight: bold;'> " + appointmentDate + " a las " + appointmentHour + " </span><br/><a style='color: #54ace2;font-weight: bold;font-size: 20px;' href='" + frontEndUrl + "/confirm-appointment?tokenizer=" + token + "&appointment= " + appointment + "' >Confirmar</a><a style='color: red;font-weight: bold;font-size: 20px;margin-left:5px' href='" + frontEndUrl + "/cencel-appointment?tokenizer=" + token + "&appointment= " + appointment + "  ' >Cancelar</a></div></div> "

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "¡Tienes una cita pendiente!")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendEmailConfirmationToPatient(mail string, phone string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = "<div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: #0f76b0;font-size: 20px;font-weight: bold;'>Se ha confirmado tu cita</span><br/>No dudes en llamar al " + phone + "</div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "¡Tienes una cita pendiente!")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}

func sendEmailCancelationToPatient(mail string, phone string) {

	var config = Config.Config{}
	config.Read()

	fmt.Println("Trying to send email! " + mail)

	var htmlContentMessage = "<div style='width:100%;text-align:center'><div><img src='" + logo + "' alt='logoclic-02' border='0'></div><br><div>	<span style='color: red;font-size: 20px;font-weight: bold;'>Se ha cancelado tu cita</span><br/>Confirma que paso al " + phone + "</div></div>"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.Email)

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "¡Tienes una cita pendiente!")

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", htmlContentMessage)

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailSmtp, 587, config.Email, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//panic(err)
	}
	fmt.Println("Email Sent Successfully!")

	return
}
