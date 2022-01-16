package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

type get struct {
	id int
	imagelinks []string
	colors []string
	cat string
	count int
}
const connStr = "user=postgres password=991155 dbname=wildberries sslmode=disable host=db"

func parseCount() string {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, _ := db.Query("Select id, imagelinks, category, count, colors from items")
	info:=""
	count := 0
	for res.Next() {
		count++
		c := get{}
		res.Scan(&c.id, &c.imagelinks, &c.cat, &c.count, &c.colors)
		x := strconv.Itoa(c.id)
		xx := strconv.Itoa(c.count)
		y := strconv.Itoa(len(c.imagelinks))
		yy := strconv.Itoa(len(c.colors))
		info += "ID: " + x + " Category: " + c.cat + " Items count: " + xx + " Images count: " + y + " Colors count: " + yy + "\n"
	}
	countS := strconv.Itoa(count)
	return "COUNT: " + countS + "\n" +  info
}
func main() {
	time.Sleep(time.Second * 20)
	fmt.Println("YES !!")
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, parseCount())
	})

	fmt.Println("Server listening!")
	http.ListenAndServe(":80", r)
}
