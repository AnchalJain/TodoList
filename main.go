package main

import (

	"database/sql"
	"log"
	"net/http"
	"text/template"
	_"github.com/go-sql-driver/mysql"
)

type Todo struct {
	Id int
	Title string
	Description string
}

func dbConn() (db *sql.DB) {
	dbDriver:="mysql"
	dbUser:="root"
	dbPass:="forceofnature"
	dbName:="goblog"

	db, err:=sql.Open(dbDriver,dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl=template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db:=dbConn()
	selDb, err:=db.Query("SELECT * FROM todo ORDER BY id DESC")

	if err!= nil {
		panic(err.Error())
	}

	todo:=Todo{}
	res:=[]Todo{}
	for selDb.Next() {
		var id int
		var title, description string
		err = selDb.Scan(&id, &title, &description)
		if err != nil {
			panic(err.Error())
		}
		todo.Id=id
		todo.Title=title
		todo.Description=description
		res=append(res, todo)
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

//func Show(w http.ResponseWriter, r *http.Request) {
//	db:=dbConn()
//	nId:=r.URL.Query().Get("id")
//	log.Println(nId)
//	selDb, err:= db.Query("SELECT * FROM todo WHERE id=?", nId)
//	if err!=nil {
//		panic(err.Error())
//	}
//	todo:=Todo{}
//	res:=[]Todo{}
//	for selDb.Next() {
//		var id int
//		var title, description string
//		err=selDb.Scan(&id, &title, &description)
//		if err != nil {
//			panic(err.Error())
//		}
//		todo.Id=id
//		todo.Title=title
//		todo.Description=description
//		res=append(res, todo)
//	}
//
//	w.Header().Set("Content-Type", "text/html")
//	tmpl.ExecuteTemplate(w, "Show", res)
//	defer db.Close()
//}

func New(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db:=dbConn()
	nId:=r.URL.Query().Get("id")
	selDb, err:=db.Query("SELECT * FROM todo WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	todo:=Todo{}
	for selDb.Next() {
		var id int
		var title, description string
		err=selDb.Scan(&id, &title, &description)
		if err != nil {
			panic(err.Error())
		}
		todo.Id=id
		todo.Title=title
		todo.Description=description
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "Edit", todo)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db:=dbConn()
	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")
		insForm, err := db.Prepare("INSERT INTO todo(title, description) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(title, description)
//		log.Println("INSERT: Title: "+ title + "Description: "+ description)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
	title:=r.FormValue("title")
	description:=r.FormValue("description")
	id:=r.FormValue("uid")
	insForm, err := db.Prepare("UPDATE todo SET title=?, description=? WHERE id=?")
	if err != nil {
		panic (err.Error())
	}
	insForm.Exec(title, description, id)
	log.Println("UPDATE: Title: "+ title + " | Description: " + description)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    todo:=r.URL.Query().Get("id")
    log.Println(todo)
    delForm, err := db.Prepare("DELETE FROM todo WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    log.Println("DELETE")
    delForm.Exec(todo)
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on : http://localhost:8080")
	http.HandleFunc("/", Index)
//	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
