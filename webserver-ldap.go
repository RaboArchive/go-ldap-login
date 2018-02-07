package main

import (
	"fmt"
	ldap "go-ldap-client"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Global Var
var client *ldap.LDAPClient

func authUser(login string, pass string) (bool, error) {
	// It is the responsibility of the caller to close the connection

	ok, user, err := client.Authenticate(login, pass)
	if err != nil {
		fmt.Printf("Error authenticating user %s: %+v", login, err)
	}
	if !ok {
		fmt.Printf("Authenticating failed for user %s", login)
	}
	log.Printf("User: %+v", user)
	return ok, err
}

/*
 * Web server
 * Routes
 */
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parsing parameters
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "There is nothing to see here John Snow, go there : http://localhost:9090/login") // send data to client side
}

func login(w http.ResponseWriter, r *http.Request) {
	html, err := loadPage("login")
	if err != nil {
		return
	}
	fmt.Fprintf(w, string(html))
}

func auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parsing parameters
	login := strings.Join(r.Form["login"], "")
	pass := strings.Join(r.Form["pass"], "")
	fmt.Printf("%s : %s\n", login, pass)
	res, err := authUser(login, pass)
	var response string
	if err != nil {
		fmt.Printf("Failure\n")
		response = "Unable to log the user : " + login
	}
	if res {
		fmt.Printf("Success\n")
		response = login + "User successfuly connected"
	}
	fmt.Fprintf(w, response)
}

// Utils
func loadPage(title string) ([]byte, error) {
	filename := "pages/" + title + ".html"
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", html)
	return html, nil
}

func main() {
	// Connect to LDAP
	// We use this to acces to the iut LDAP
	//ssh -L 3890:ldap.iut2.upmf-grenoble.fr:389 prudhomj@transit.iut2.upmf-grenoble.fr
	client = &ldap.LDAPClient{
		Base:         "ou=people,dc=iut2,dc=upmf-grenoble,dc=fr",
		Host:         "localhost",
		Port:         3890,
		UseSSL:       false,
		BindDN:       "",
		BindPassword: "",
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	defer client.Close()

	// Router
	http.HandleFunc("/", index) // set router
	http.HandleFunc("/login", login)
	http.HandleFunc("/auth", auth)
	// Starting the server
	err1 := http.ListenAndServe(":9090", nil) // set listen port
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err1)
	}
}
