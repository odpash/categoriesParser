package main

import (
	"database/sql"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const categoryFilename = "category.json"
const connStr = "user=postgres password=991155 dbname=wildberries sslmode=disable host=db"
//const connStr = "postgres://postgres:991155@0.0.0.0:5433/wildberries?sslmode=disable"

type Categories struct {
	Categories []Category
}

type Category struct {
	Name    string
	PageUrl string
}

func WriteDb(info Categories) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec("delete from categories *")
	for _, value := range info.Categories {
		_, e := db.Exec("insert into categories (name, link) values ($1, $2)",
			value.Name, value.PageUrl)
		if e != nil {
			fmt.Println(e)
		}
	}
	fmt.Println("All is OK!")
}

func scrapCategoriesCycle(c []byte, newCategories Categories) Categories {
	_, err := jsonparser.ArrayEach(c, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		cNew, _, _, childError := jsonparser.Get(value, "childs")
		if childError != nil {
			name, _, _, _ := jsonparser.Get(value, "name")
			pageUrl, _, _, _ := jsonparser.Get(value, "pageUrl")
			if strings.Contains(string(pageUrl), "catalog") {
				if !strings.Contains(string(pageUrl), "https://digital") {
					newCategory := Category{
						Name:    string(name),
						PageUrl: "https://www.wildberries.ru" + string(pageUrl),
					}
					isIn := false
					for _, v := range newCategories.Categories {
						if v.Name == newCategory.Name {
							isIn = true
						}
					}
					if !isIn {
						newCategories.Categories = append(newCategories.Categories, newCategory)
					}
				}
			}
		} else {
			newCategories = scrapCategoriesCycle(cNew, newCategories)
		}

	})

	if err != nil {
		return newCategories
	} else {
		return newCategories
	}
}

func scrapCategories() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f20597c3014e4699969af0244a66a6f8@o1108001.ingest.sentry.io/6135375",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	start := time.Now()
	var newCategories Categories
	url := "https://www.wildberries.ru/gettopmenuinner?lang=ru"
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	c, _, _, _ := jsonparser.Get(body, "value", "menu")
	newCategories = scrapCategoriesCycle(c, newCategories)
	WriteDb(newCategories)
	sentry.CaptureMessage("[1/4] Парсер категорий выполнил задачу за " + time.Since(start).String() + " и получил " + strconv.Itoa(len(newCategories.Categories)) + " категорий.")
}

func main() {
	time.Sleep(time.Second * 10)
	fmt.Println("Category Started")
	for {
		scrapCategories()
		time.Sleep(time.Hour)
	}
}