package main

import(
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type User struct{
	gorm.Model
	Name string `json:"Name"`
	Food string `json:"Food"`
	Sport string `json:"Sport"`
}

func getConnection() (*gorm.DB){

	dsn := "root:root@(127.0.0.1:3308)/gorm_sql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err!=nil{
		panic(err.Error)
	}
	fmt.Println("connected")

	return db;

}

func createTable(){

	db:=getConnection()
	db.Debug().AutoMigrate(&User{})
	fmt.Println("created")
}

func getAllUsers(w http.ResponseWriter,r *http.Request){

	db:=getConnection()
	var u []User
	result:=db.Find(&u)
	fmt.Println(result.RowsAffected)
	fmt.Println(u)

	data,_:=json.Marshal(u)
	w.Write(data)
}

func getOneUser(w http.ResponseWriter,r *http.Request){

	db:=getConnection()
	var u User
	id:=mux.Vars(r)["id"]
	result:=db.First(&u,id)
	fmt.Println(result.RowsAffected)

	data,_:=json.Marshal(u)
	w.Write(data)
}

func insertUser(w http.ResponseWriter,r *http.Request){

	db:=getConnection()

	var u User
	data,_:=ioutil.ReadAll(r.Body)
	_=json.Unmarshal(data,&u)

	result:=db.Create(&u)
	if result.Error!=nil{
		panic(result.Error)
	}

	data,_=json.Marshal(u)
	w.Write(data)

}

func updateUser(w http.ResponseWriter,r *http.Request){

	db:=getConnection()
	var u User
	id:=mux.Vars(r)["id"]

	_=db.First(&u,id)

	data,_:=ioutil.ReadAll(r.Body)
	_=json.Unmarshal(data,&u)

	db.Save(&u)

	data,_=json.Marshal(u)
	w.Write(data)
}

func deleteUser(w http.ResponseWriter,r *http.Request){
	db:=getConnection()
	id:=mux.Vars(r)["id"]
	db.Delete(&User{},id)
	response,_:=json.Marshal(struct{Msg string}{"deleted successfully"})
	w.Write(response)
}

func main(){

	fmt.Println("listening on port 8080...")
	r:=mux.NewRouter()
	r.HandleFunc("/getAllUsers",getAllUsers).Methods("GET")
	r.HandleFunc("/getOneUser/{id}",getOneUser).Methods("GET")
	r.HandleFunc("/insertUser",insertUser).Methods("POST")
	r.HandleFunc("/updateUser/{id}",updateUser).Methods("PUT")
	r.HandleFunc("/deleteUser/{id}",deleteUser).Methods("DELETE")
	http.ListenAndServe(":8080",r)

}