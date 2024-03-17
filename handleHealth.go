package main

import "net/http"

func checkHealth(w http.ResponseWriter, r *http.Request) {
	type respStruct struct {
		Msg string `json:"msg"`
	}

	respObj := respStruct{
		Msg: "OK",
	}

	respWithJson(w, 200, respObj)
}

func checkError(w http.ResponseWriter, r *http.Request) {
	respWithError(w, 200, "Error OK")
}
