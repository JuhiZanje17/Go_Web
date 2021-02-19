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
			PostID    int       `json:"post_id"`
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			Body      string    `json:"body"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`	
}

func main(){

	url:="https://gorest.co.in/public-api/comments"

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

	file_content,_:=json.MarshalIndent(data,""," ")
	err=ioutil.WriteFile("comment.json",file_content,0644)
	if err!=nil{
		log.Fatalln(err)
	}


}