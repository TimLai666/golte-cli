package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := ginRouter()

	fmt.Println("Serving on :8000")
	http.ListenAndServe(":8000", r)
}
