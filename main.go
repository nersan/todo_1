package main

import (
	"fmt"
	"net/http"
	// "os"
)

func main() {
	// err := http.ListenAndServe(
	// 	":18080",
	// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintf(w, "Hello,%s!", r.URL.Path[1:])
	// 	}),
	// )
	// if err != nil {
	// 	fmt.Fprintln(os.Stdout, "Error: ", err)
	// 	os.Exit(1)
	// }
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	s.ListenAndServe()
}
