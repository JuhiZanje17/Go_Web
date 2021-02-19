package main

import (
	"database/sql"
	_"go-sql-driver/mysql"
	"fmt"
)


func main(){
	db,err:=sql.Open("mysql","root:@(127.0.0.1:3308)/go_db?parseTime=true")
	//fmt.Printf("db=%T",db)
	if err!=nil{
		panic(err.Error())
	}
	err=db.Ping()
	if err!=nil{
		panic(err.Error())
	}

	//insert_table(db)

}

type user struct{
	id int
	name string
	food string
	sport string
}

func create_table(db *sql.DB){

	query:=`create table users(id int auto_increment,name varchar(20) not null,food varchar(10) not null,sport varchar(20) not null,primary key(id))`
	_,err:=db.Exec(query)
	if err!=nil{
		panic(err.Error())
	}	
}

func select_row(db *sql.DB){

	query:=`select * from users where id=?`
	fmt.Println("Enter id to search record:")
	var id int
	var name,food,sport string
	fmt.Scan(&id)
	err:=db.QueryRow(query,id).Scan(&id,&name,&food,&sport)
	if err!=nil{
		panic(err.Error())
	}
	fmt.Printf("name=%v food=%v sport=%v",name,food,sport)
}

func select_table(db *sql.DB){

	query:=`select * from users`
	rows,err:=db.Query(query)
	if err!=nil{
		panic(err.Error())
	}
	var users[] user
	for rows.Next(){
		var u user
		err=rows.Scan(&u.id,&u.name,&u.food,&u.sport)
		if err!=nil{
			panic(err.Error())
		}
		users=append(users,u)
	}

	fmt.Println(users)
}

func insert_table(db *sql.DB){

	insert_query:=`insert into users values (null,?,?,?)`

	result,err:=db.Exec(insert_query,"jhanvi","panipuri","secret")
	if err!=nil{
		panic(err.Error())
	}
	fmt.Println(result.LastInsertId())
}

func update_table(db *sql.DB){

	query:=`update users set food=? where id=?`
	var id int
	var food string
	fmt.Println("Enter id to update record:")
	fmt.Scan(&id)
	fmt.Println("Enter food to update record:")
	fmt.Scan(&food)

	result,err:=db.Exec(query,food,id)
	if err!=nil{
		panic(err.Error())
	}

	fmt.Println(result.RowsAffected())
}

func delete_table(db *sql.DB){

	query:=`delete from users where id=?`
	var id int
	fmt.Println("Enter id to delete record:")
	fmt.Scan(&id)

	result,err:=db.Exec(query,id)
	if err!=nil{
		panic(err.Error())
	}

	fmt.Println(result.RowsAffected())
}


