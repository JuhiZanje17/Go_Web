package main

import (

	 "fmt"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	// "io/ioutil"
	// "database/sql"
	// _ "github.com/lib/pq"
)

func main(){

	fmt.Println("listening on port 8080...")

	r := mux.NewRouter().StrictSlash(true)

	// Choose the folder to serve
	staticDir := "/static/"

	// Create the route
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	r.HandleFunc("/upload",getFile).Methods("GET")	
	http.ListenAndServe(":8000",r)
}

func getFile(w http.ResponseWriter,r * http.Request){

	fmt.Println("In file")
	//http.ServeFile(w, r, "upload.html")
	tmpl,_:=template.ParseFiles("upload.html")
	err:=tmpl.Execute(w,nil)
	if err!=nil{
		panic(err.Error)
	}
}