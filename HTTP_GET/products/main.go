package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"log"
	"time"
)

type struct_to_generate struct{
	
	Code int `json:"code"`
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
			Pages int `json:"pages"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
		
}

func main(){

	url:="https://gorest.co.in/public-api/products"

	resp,err:=http.Get(url)

	if err!=nil{
		log.Fatalln(err)
	}

	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		log.Fatalln(err)
	}

	var data struct_to_generate
	_=json.Unmarshal(body,&data)

	fmt.Println(data.Meta)

	file_content,_:=json.MarshalIndent(data,"","	")
	err=ioutil.WriteFile("product.json",file_content,0644)
	if err!=nil{
		log.Fatalln(err)
	}


}