package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Send(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println(err)

		return
	}
}

func Fail(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.WriteHeader(code)

	err2 := json.NewEncoder(w).Encode(err.Error())
	if err != nil && err2 != nil {
		fmt.Println(err)
	}
}
