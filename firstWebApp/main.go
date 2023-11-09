package main

import (
	"net/http"
	"app/handlers"
)


func main() {
	defer handlers.Db.Close()
	
	http.HandleFunc("/delete/", handlers.Delete)
	http.HandleFunc("/new/", handlers.New)
	http.HandleFunc("/home/", handlers.Home)

	http.Handle("/", http.FileServer(http.Dir("./public/")))
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}
}
