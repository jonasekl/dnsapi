package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/googollee/go-rest"
	"net/http"
	"time"
)

type DNSAPI struct {
	rest.Service `prefix:"/dnsapi" mime:"application/json" charset:"utf-8"`

	CreateEntry rest.Processor `method:"POST" path:"/create"`

	post  map[string]string
	watch map[string]chan string
}

type EntryArg struct {
	//	Id int `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}



// Post like:
// > curl "/dnsapi/create" -d '{"name":"mymachine", "content":"10.1.2.3"}'
//
// No response
func (r DNSAPI) HandleCreateEntry(arg EntryArg) {
	db, err := sql.Open("mysql", "myuser:mypass@/mydb")

	fmt.Println("create entry:" + arg.Name + ". content:" + arg.Content)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println("Will now attempt to delete other " + arg.Name)
	stmt, _ := db.Prepare("DELETE FROM records where name=?")
	res, err := stmt.Exec(arg.Name)
	if err != nil{
		panic(err.Error())
	}
	fmt.Println(res)

	stmt, _ = db.Prepare("INSERT INTO records(domain_id, name, type, content, ttl, prio, change_date) VALUES(?, ?, ?, ?, ?, ?, ?)")
	res, err = stmt.Exec(1, arg.Name, "A", arg.Content, 3600, 0, time.Now().Unix())
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)

}


func main() {
	handler, err := rest.New(&DNSAPI{
		post:  make(map[string]string),
		watch: make(map[string]chan string),
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Starting tha webservice")
	http.ListenAndServe("0.0.0.0:8889", handler)
}
