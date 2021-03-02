/*
go get -u gorm.io/driver/mysql
go get -u gorm.io/driver/postgres
*/
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"fmt"
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

func insertUser(){

	db:=getConnection()
	u:=User{Name:"Ju",Food:"kuch bhi",Sport:"nothing"}	
	result:=db.Create(&u)
	if result.Error!=nil{
		panic(result.Error)
	}

	fmt.Println("user id=",u.ID)
	fmt.Println("rows affected=",result.RowsAffected)

}

func getAllUsers(){

	db:=getConnection()
	var u []User
	result:=db.Find(&u)
	fmt.Println(result.RowsAffected)
	fmt.Println(u)
}

func getOneUser(){

	db:=getConnection()
	var u User
	var id int
	fmt.Println("enter id=")
	fmt.Scan(&id)
	result:=db.First(&u,id)
	fmt.Println(result.RowsAffected)
	fmt.Println(u)
}

func updateUser(){

	db:=getConnection()
	var u User
	var id int
	fmt.Println("enter id to update=")
	fmt.Scan(&id)
	_=db.First(&u,id)
	u.Name="u1"
	u.Food="u2"
	u.Sport="u3"
	db.Save(&u)
}

func deleteUser(){
	db:=getConnection()
	var id int
	fmt.Println("enter id to delete=")
	fmt.Scan(&id)
	db.Delete(&User{},id)
}


func main() {

	//createTable()
	// insertUser()
	// getAllUsers()
	// getOneUser()
	//updateUser()
	deleteUser()
	getAllUsers()

}

