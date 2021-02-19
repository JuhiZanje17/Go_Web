package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"database/sql"
	_ "github.com/lib/pq"
)

type Image struct{
	Id int
	Name string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "go_postgre"
)

func dbConn() (db *sql.DB){

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
	panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func ctreateTable(){
	//CREATE
	db:=dbConn()
	defer db.Close()
	query:=`create table images_go(id serial primary key,name varchar(30))`
	_,err:=db.Exec(query)
	if err!=nil{
		panic(err.Error())
	}	

	fmt.Println("Successfully created!")
}

func getFile(w http.ResponseWriter,r * http.Request){

	tmpl,_:=template.ParseFiles("upload.html")
	err:=tmpl.Execute(w,nil)
	if err!=nil{
		panic(err.Error)
	}
}

func uploadFile(w http.ResponseWriter,r * http.Request){

	//parses a request body as multipart/form-data.The whole request body is parsed and up to a total of maxMemory bytes	
	r.ParseMultipartForm(10 << 20)
	file,handler,err:=r.FormFile("file")
	if err!=nil{
		panic(err.Error)
	}
	fmt.Println("File name=",handler.Filename)
	fmt.Println("File size=",handler.Size)
	fmt.Println("File header=",handler.Header)

	defer file.Close()

	Path := filepath.Join("static/", handler.Filename)
		fmt.Println(Path)
		image, err := os.Create(Path)
		if err != nil {
			panic(err)
		}

		defer image.Close()
		io.Copy(image, file)

	db:=dbConn()
	defer db.Close()

	query:=`insert into images_go (name) values ($1)`
	_,err=db.Exec(query,handler.Filename)
	if err!=nil{
		panic(err.Error())
	}

	fmt.Println("Successfully inserted!")
	http.Redirect(w, r, "/show", 301)

/*
	using tempfile:
	tempFile,err:=ioutil.TempFile("static","upload-*.png")
	if err!=nil{
		panic(err.Error)
	}
	defer tempFile.Close()
	fmt.Printf("%T",tempFile.Name())

	dataBytes,err:=ioutil.ReadAll(file)
	if err!=nil{
		panic(err.Error)
	}
	tempFile.Write(dataBytes)
	w.Write([]byte("uploaded successfully!!"))
*/
}

func showImages(w http.ResponseWriter,r * http.Request){

	tmpl,_:=template.ParseFiles("display.html")
	err:=tmpl.Execute(w,nil)
	if err!=nil{
		panic(err.Error)
	}

	var images []Image
	db := dbConn()
	query := "SELECT *FROM images_go"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("selected successfully")
	for rows.Next() {
		var image Image
		err = rows.Scan(&image.Id, &image.Name)
		if err != nil {
			panic(err)
		}
		images = append(images, image)
	}
	fmt.Println(images)
	tmpl.ExecuteTemplate(w, "display.html", images)
}

func main(){

	fmt.Println("listening on port 8080...")
	//ctreateTable()

	r := mux.NewRouter().StrictSlash(true)
	// Choose the folder to serve
	staticDir := "/static/"
	// Create the route
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	r.HandleFunc("/upload",getFile).Methods("GET")	
	r.HandleFunc("/upload",uploadFile).Methods("POST")
	r.HandleFunc("/show",showImages).Methods("GET")
	http.ListenAndServe(":8080",r)
}