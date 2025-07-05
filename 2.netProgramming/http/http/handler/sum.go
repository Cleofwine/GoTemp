package handler

import (
	"fmt"
	"net/http"
	"strconv"
)

func SumHandler(w http.ResponseWriter, req *http.Request) {
	a := req.FormValue("a")
	b := req.FormValue("b")
	intA, _ := strconv.Atoi(a)
	intB, _ := strconv.Atoi(b)
	fmt.Fprintf(w, "a + b = %d", intA+intB)
}
