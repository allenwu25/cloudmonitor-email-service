package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)


func readFromTemplate(templatename string) (string) {
	body, templateerr := ioutil.ReadFile(templatename)
	if (templateerr != nil) {
		fmt.Print(templateerr)
	}

	bodycontent := string(body)
	return bodycontent
}


type test_struct struct {
	Emails []string
}



func sendMail(Emails [] string) {
	from := mail.NewEmail("Cloudmonitor", "a52wu@uwaterloo.ca")
	htmlContent := readFromTemplate("template.html")
	content := mail.NewContent("text/html", htmlContent)

	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.AddContent(content)

	personalization := mail.NewPersonalization()

	for _, email := range Emails {
		to := mail.NewEmail("Cloudmonitor user", email)
		personalization.AddTos(to)
	}

	personalization.Subject = "API Test Error Alert from Cloudmonitor"

	m.AddPersonalizations(personalization)
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
	  log.Println(err)
	} else {
	  fmt.Println(response.StatusCode)
	  fmt.Println(response.Body)
	  fmt.Println(response.Headers)
	}
}



func parseJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t test_struct
	err := decoder.Decode(&t)

	if err != nil {
		http.Error(w, "Incorrect parameters passed", http.StatusBadRequest)
	} else {
		sendMail(t.Emails)
		fmt.Fprint(w, "Successfully Sent Emails")
	}

}


func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Email Server Running")
}


func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/sendemail", parseJSON).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	handleRequests()
}
