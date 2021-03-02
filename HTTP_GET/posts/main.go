package main
import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
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
			Body      string    `json:"body"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`
}

func main(){
	url:="https://gorest.co.in/public-api/posts"

	resp,err:=http.Get(url)

	if err!=nil{
		log.Fatalln(err)
	}

	//reading the body of resp (url)
	body,err:=ioutil.ReadAll(resp.Body)

	//converting read byte(json) data into struct (unmarshaling)
	var data struct_to_generate

	err=json.Unmarshal(body,&data)

	if err!=nil{
		log.Fatalln(err)
	}

	fmt.Println(data.Meta)

	//generating a json file having the content of the url
	file_content,_:=json.MarshalIndent(data,"","	")
	err=ioutil.WriteFile("post.json",file_content,0644)
	if err!=nil{
		log.Fatalln(err)
	}
}