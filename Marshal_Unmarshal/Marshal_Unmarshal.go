/*As with all structs in Go, it’s important to remember
that only fields with a capital first letter are visible to external programs like the JSON Marshaller*/

package main

import (
	"fmt"
	"encoding/json"
)

func main(){

	type add struct{
		Area string
		Country string
	}

	type pDet struct{
		Id int
		Name string
		Address add
	}

	s1:=[] pDet{
		{1,"Juhi",add {"Jivaraj","IND"}},
		{2,"Jhanvi",add {"Jivaraj","UK"}},
	}

	emp, _ := json.Marshal(s1[0])
	fmt.Println(s1[0])
	fmt.Println(string(emp))

	str_emp:=string(emp)
	
	type unmarshal_struct struct{
		Id int `json:"id"`
		Name string `json:"name"`
		Address add `json:"address"`
	}

	bytes:=[]byte(str_emp)
	var ans unmarshal_struct
	json.Unmarshal(bytes,&ans)
	fmt.Println(ans)
}



