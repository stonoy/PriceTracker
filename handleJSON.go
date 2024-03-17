package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respWithError(w http.ResponseWriter, code int, msg string) {
	type respStruct struct {
		Error string `json:"error"`
	}

	respObj := respStruct{
		Error: msg,
	}

	if code > 499 {
		log.Printf("Error in Server Code: %v, Msg: %v/n", code, msg)
	}

	respWithJson(w, code, respObj)
}

func respWithJson(w http.ResponseWriter, code int, respObj interface{}) {
	dat, err := json.Marshal(respObj)
	if err != nil {
		respWithError(w, 500, "can not marshal respObj")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
