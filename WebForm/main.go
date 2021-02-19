package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"os"
	"log"
)

type info struct{
	Name string
	Food string
	Sport string
}

func main(){	

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl,_:=template.ParseFiles("form.html")

	http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){		
		
		if r.Method==http.MethodGet {
			
			tmpl.Execute(w,nil)
			return
		}		
		
		data:=info{
			r.FormValue("name"),
			r.FormValue("food"),
			r.FormValue("sport"),
		}

		fmt.Println(data)

		file_content,_:=json.MarshalIndent(data,"","	")
		
		f,err := os.OpenFile("Information.json",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
		
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(file_content)
		if err != nil {
			f.Close() 
			log.Fatal(err)
		}
		f.Close()

		tmpl.Execute(w,struct{Success bool}{true})

	})

	fmt.Println("listening")

	http.ListenAndServe(":8080",nil)	
}