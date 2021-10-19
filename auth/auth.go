package auth

import (
	"log"
	"net/http"
)

func IsAuthenticated(r *http.Request) bool {
	//static test elements for basic authentiation
	e := "hart87@gmail.com"
	p := "password"

	em, pass, ok := r.BasicAuth()

	if !ok {
		log.Print("Not Okay")
		return false
	}

	if e != em {
		log.Print("Bad Email")
		return false
	}

	if p != pass {
		log.Print("Bad Password")
		return false
	}

	return true
}
