package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	gomail "gopkg.in/mail.v2"
)

var smtphost *string = flag.String("smtphost", "", "smtp host, e.g. smtp.example.com")
var smtpport *int = flag.Int("smtpport", 587, "the port to connect to the smtp")
var smtpuser *string = flag.String("smtpuser", "", "username for the smtp")
var smtppassword *string = flag.String("smtppassword", "", "password for the smtp user")
var smtpoverride *bool = flag.Bool("smtpoverride", true, "true - allows to pass smtp parameters in the json call, false will always use the config smtp data")

type realAttachments map[string]string

type Email struct {
	// email parameters
	From        string          `json:"from"`
	To          []string        `json:"to"`
	Cc          []string        `json:"cc"`
	Bcc         []string        `json:"bcc"`
	Subject     string          `json:"subject"`
	Message     string          `json:"message"`
	Attachments realAttachments `json:"attachments"`

	//smtp parameters
	Smtphost     *string `json:"smtphost"`
	Smtpport     *int    `json:"smtpport"`
	Smtpuser     *string `json:"smtpuser"`
	Smtppassword *string `json:"smtppassword"`
}

func main() {
	flag.Usage = func() {
		fmt.Println("json2smtp utility https://www.c2kb.com/json2smtp v1.0.1 2023-11-13")
		fmt.Println("Get json input and calls smtp - function as a json to smtp proxy")
		fmt.Println("Options:")
		flag.PrintDefaults()

		// Custom help information.
		fmt.Println("\nExample:")
		fmt.Printf("\t%v --port=8200 --smtphost='smtp.example.com' --smtpport=587 --smtpuser='username' --smtppassword='password' --smtpoverride=false", os.Args[0])
		fmt.Println("\nRun in the background:")
		fmt.Printf("\tnohup %v --port=8200 --smtphost='smtp.example.com' --smtpport=587 --smtpuser='username' --smtppassword='password' >> logfile.log 2>&1 &", os.Args[0])
		fmt.Println("\n\nParametrs to pass in the json call:")
		fmt.Println(`{
	"from": "john doe <john@example.com>",
	"to": ["kermit@muppets.com", "oneperson@example.com"],
	"cc": ["email1@example.com"],
	"bcc": ["secret@example.com"],
	"subject": "email subject line",
	"message": "message body in text/html to be sent",
	"attachments": {"filename.pdf": "base64 file encoded", "anotherfilename.txt": "base64 file encoded"},
	"smtphost": "smtp.example.com - optional parameter",
	"smtpport": 587 - optional paramater,
	"smtpuser": "username - optional parameter",
	"smtppassword": "password - optional parameter"
}`)

		fmt.Println("\nRecommendation: Put the service behind caddy for SSL/TLS connection encryption")
	}

	port := flag.Int("port", 8080, "the port to listen on")

	flag.Parse()

	log.Println("json2smtp server started, listening on port: ", *port, " host:", *smtphost, " allow json smtp information:", *smtpoverride)

	http.HandleFunc("/", handlejson2smtp)

	http.ListenAndServe(":"+fmt.Sprint(*port), nil)

	log.Println("json2smtp server ended")
}

func handlejson2smtp(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into an Email object.
	var email Email
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Error decoding JSON request body: %v"}`, err)
		log.Printf("error decoding JSON request body: %v", err)
		return
	}

	// Build the smtp data
	localsmtphost := email.Smtphost
	if localsmtphost == nil || !*smtpoverride {
		localsmtphost = smtphost
	}

	localsmtpport := email.Smtpport
	if localsmtpport == nil || !*smtpoverride {
		localsmtpport = smtpport
	}

	localsmtpuser := email.Smtpuser
	if localsmtpuser == nil || !*smtpoverride {
		localsmtpuser = smtpuser
	}

	localsmtppassword := email.Smtppassword
	if localsmtppassword == nil || !*smtpoverride {
		localsmtppassword = smtppassword
	}

	// sanity check we have smtp information
	if localsmtphost == nil || localsmtpport == nil || localsmtpuser == nil || localsmtppassword == nil {
		fmt.Fprintf(w, `{"error": "Missing smtp host data for sending"}`)
		log.Printf("Missing smtp host data for sending")
		return
	}

	// Decode attachments
	ra := make(realAttachments, 10)
	for attachmentName, data := range email.Attachments {
		decodedAttachment, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			fmt.Fprintf(w, `{"error": "Error decoding base64 attachment: %v"}`, err)
			log.Printf("Error decoding base64 attachment: %v", err)
			return
		}

		ra[attachmentName] = string(decodedAttachment)
	}

	// Send the email. (TODO: make the calls to use go routines for parallel email sending)
	err = sendEmail(*localsmtphost, *localsmtpport, *localsmtpuser, *localsmtppassword, email.From, email.To, email.Cc, email.Bcc, email.Subject, email.Message, ra)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Error sending email: %v"}`, err)
		log.Printf("Error sending email: %v", err)
		return
	}

	// Respond to the client with a success message.
	fmt.Fprintf(w, `{"success": true, "to": "%v", "subject": "%v"}`, email.To, email.Subject)
}

func sendEmail(smtphost string, smtpport int, smtpuser, smtppassword, from string, to, cc, bcc []string, subject, message string, ra realAttachments) error {
	// clean the arrays
	to = deleteEmpty(to)
	cc = deleteEmpty(cc)
	bcc = deleteEmpty(bcc)

	log.Printf("Sending email from: %v, to: %v, cc: %v, bcc: %v, subject: %v", from, to, cc, bcc, subject)
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", to...)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("text/plain", message)
	m.SetBody("text/html", message)

	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}
	if len(bcc) > 0 {
		m.SetHeader("Bcc", bcc...)
	}

	// Add the attachments
	for name, data := range ra {
		log.Println("attachment name:", name, " data length:", len(data))
		reader := strings.NewReader(data)
		//defer reader.Close()
		m.AttachReader(name, reader)
	}

	// Settings for SMTP server
	d := gomail.NewDialer(smtphost, smtpport, smtpuser, smtppassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	log.Printf("Email success from: %v, to: %v, subject: %v", from, to, subject)
	return nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
