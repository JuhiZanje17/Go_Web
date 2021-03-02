
package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
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
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			Gender    string    `json:"gender"`
			Status    string    `json:"status"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`	
}

func main(){

	fmt.Println("USERS:")
	url:="https://gorest.co.in/public-api/users"

	resp,err:=http.Get(url)

	if err!=nil {
		fmt.Println(err)
	}	
	
	body,err:=ioutil.ReadAll(resp.Body)
	//body is []byte

	if err!=nil {
		fmt.Println(err)
	}

	var data struct_to_generate

	err=json.Unmarshal(body,&data)

	if err!=nil {
		fmt.Println(err)
	}

	fmt.Println(data.Meta)

	// mar,_:=json.Marshal(data)
	// fmt.Println("marshal\n",string(mar))

	file_content,_:=json.MarshalIndent(data,"","	")
	err=ioutil.WriteFile("user.json",file_content,0644)

	if err!=nil{
		fmt.Println(err)
	}
	




}