package main

import (
	"database/sql"
	_"go-sql-driver/mysql"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
)

func dbConn() (db *sql.DB) {

	db,err:=sql.Open("mysql","root:root@(127.0.0.1:3308)/go_db?parseTime=true")
    if err != nil {
        panic(err.Error())
    }

	err=db.Ping()
	if err!=nil{
		panic(err.Error())
	}
    return db
}

type user struct{
	Id int `json:"ID"`
	Name string `json:"Name"`
	Food string `json:"Food"`
	Sport string `json:"Sport"`
}

func allUsers(w http.ResponseWriter,r *http.Request){
	db:=dbConn()
	query:=`select * from users`
	rows,err:=db.Query(query)
	if err!=nil{
		panic(err.Error())
	}

	//store retrvied data in an array of struct user
	var users []user
	for rows.Next(){
		var u user
		err=rows.Scan(&u.Id,&u.Name,&u.Food,&u.Sport)
		if err!=nil{
			panic(err.Error())
		}
		users=append(users,u)
	}

	data,_:=json.Marshal(users)
	w.Write(data)

}

func find(w http.ResponseWriter,r *http.Request){
	
	id:=r.URL.Query().Get("id")
	db:=dbConn()
	query:=`select * from users where id=?`	
	var u user
	err:=db.QueryRow(query,id).Scan(&u.Id,&u.Name,&u.Food,&u.Sport)
	if err!=nil{
		panic(err.Error())
	}

	data,_:=json.Marshal(u)
	w.Write(data)

}

func update(w http.ResponseWriter,r *http.Request){	
	
	db:=dbConn()
	var u user
	
	u.Id,_=strconv.Atoi(r.FormValue("id"))//func ParseInt(s string, base int, bitSize int) (i int64, err error)	
	
	u.Name=r.FormValue("name")
	u.Food=r.FormValue("food")
	u.Sport=r.FormValue("sport")

	query:=`update users set name=?,food=?,sport=? where id=?`
	_,err:=db.Exec(query,u.Name,u.Food,u.Sport,u.Id)
	if err!=nil{
		panic(err.Error)
	}

	ans,_:=json.Marshal(struct{Flag bool}{false})
	w.Write(ans)
}

func delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	id := r.URL.Query().Get("Id")
	delete, err := db.Query("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()
	res, _ := json.Marshal(struct {Is_error bool}{false})
	w.Write(res)
}

func insert(w http.ResponseWriter,r *http.Request){	

	db:=dbConn()
	var u user

	data,_:=ioutil.ReadAll(r.Body)
	_=json.Unmarshal(data,&u)

	query:=`insert into users values (null,?,?,?)`
	_,err:=db.Exec(query,u.Name,u.Food,u.Sport)
	if err!=nil{
		panic(err.Error)
	}
	ans,_:=json.Marshal(struct{Flag bool}{false})
	w.Write(ans)
	
}

func main(){
	fmt.Println("hello there")
	
	http.HandleFunc("/update",update)	
	http.HandleFunc("/insert",insert)	
	http.HandleFunc("/find",find)
	 http.HandleFunc("/delete",delete)
	http.HandleFunc("/allUsers",allUsers)	
	http.ListenAndServe(":8080",nil)
}