package main

import(
	"fmt"
	"database/sql"
	_"go-sql-driver/mysql"
	"net/http"
	"html/template"
	 "strconv"
	// "strings"
)
			

type user struct{
	Id int64
	Name string    
	Food string
	Sport string
}


func main(){
	fmt.Println("hello there")
	fs:=http.FileServer(http.Dir("static/"))
		http.Handle("/static/",http.StripPrefix("/static/",fs))
	
	http.HandleFunc("/update",update_page)	
	http.HandleFunc("/insert",insert_page)	
	http.HandleFunc("/delete",delete_page)
	http.HandleFunc("/",home_page)	
	http.ListenAndServe(":8080",nil)
}
s
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

func home_page(w http.ResponseWriter,r *http.Request){		                		

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

	tmpl,_:=template.ParseFiles("home_page.html")
	err=tmpl.Execute(w,users)
 	if(err!=nil){
		panic(err.Error)
	}
}

func update_page(w http.ResponseWriter,r *http.Request){	
	
	tmpl,_:=template.ParseFiles("update_page.html")		

	if r.Method!=http.MethodPost{ 
		tmpl.Execute(w,nil)   
		return
	}     

	db:=dbConn()
	var u user
	
	u.Id,_=strconv.ParseInt(r.FormValue("id"),10,64)//func ParseInt(s string, base int, bitSize int) (i int64, err error)	
	
	u.Name=r.FormValue("name")
	u.Food=r.FormValue("food")
	u.Sport=r.FormValue("sport")

	query:=`update users set name=?,food=?,sport=? where id=?`
	result,err:=db.Exec(query,u.Name,u.Food,u.Sport,u.Id)
	if err!=nil{
		panic(err.Error)
	}
	fmt.Println(result.RowsAffected())

	http.Redirect(w, r, "/", 301)
}

func delete_page(w http.ResponseWriter,r *http.Request){

   	fmt.Println("in delete")
	db:=dbConn()
	var u user   
	query:=`delete from users where id=?`
	u.Id,_=strconv.ParseInt(r.URL.Query().Get("id"),10,64)  
	fmt.Println("uid",u.Id)
	_,err:=db.Exec(query,u.Id)
	if err!=nil{
		panic(err.Error)
	} 
	http.Redirect(w,r,"/",301)
}

func insert_page(w http.ResponseWriter,r *http.Request){	
	
	tmpl,_:=template.ParseFiles("insert_page.html")		
	if r.Method!=http.MethodPost{

		tmpl.Execute(w,nil)
		return
	}

	db:=dbConn()
	var u user
	
	u.Name=r.FormValue("name")
	u.Food=r.FormValue("food")
	u.Sport=r.FormValue("sport")

	query:=`insert into users values (null,?,?,?)`
	result,err:=db.Exec(query,u.Name,u.Food,u.Sport)
	if err!=nil{
		panic(err.Error)
	}
	fmt.Println(result.RowsAffected())
	http.Redirect(w,r,"/",301)
}

//id := r.URL.Query().Get("Id")
//funcs := template.FuncMap{"add": add}
//	temp, err := template.New("view.html").Funcs(funcs).ParseFiles("view.html")
//http.Redirect(w, r, "/", 301)