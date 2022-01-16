package main

import (
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const connStr = "user=postgres password=991155 dbname=wildberries sslmode=disable host=db"
type Id struct {
	id       int
	category string
}
type get struct {
	imagelinks []string
}

func WriteIdToPostgreSql(id int, images []string, category string) {
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	_, e := db.Exec("update items set imagelinks=$2 where id=$1", id, pq.Array(images))
	if e != nil {
		fmt.Println("Error write")
		fmt.Println(e)
	}
}

func GetDbIds() []Id {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, _ := db.Query("Select id, category from items")
	var ids []Id
	for res.Next() {
		id := Id{}
		res.Scan(&id.id, &id.category)
		ids = append(ids, id)
	}
	return ids
}

func scrapImage(id string, category string) int {
	count := 0
	var images []string
	for {
		count++
		imageLink := ""
		if len(id) == 8 {
			imageLink = "https://images.wbstatic.net/c516x688/new/" + id[0:4] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		} else if len(id) == 7 {
			imageLink = "https://images.wbstatic.net/c516x688/new/" + id[0:3] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		}

		resp, e := http.Get(imageLink)
		if e != nil {
			strId, _ := strconv.Atoi(id)
			WriteIdToPostgreSql(strId, images, category) // заменить запись
			return 1
		}
		if resp.StatusCode == 200 {
			images = append(images, imageLink)
		} else {
			strId, _ := strconv.Atoi(id)
			WriteIdToPostgreSql(strId, images, category) // заменить запись
			return 1
		}
	}
}

func scrapImages() {
	var wg sync.WaitGroup
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f20597c3014e4699969af0244a66a6f8@o1108001.ingest.sentry.io/6135375",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("[3/4] Скрипт парсера картинок запущен!")

	count := 0
	data := GetDbIds()
	for i, v := range data {
		count += 1
		wg.Add(1)
		go func(id int, category string) {
			defer wg.Done()
			scrapImage(strconv.Itoa(id), category)
		}(v.id, v.category)
		if i%50 == 0 {
			wg.Wait()
			if i%1000 == 0 {
				sentry.CaptureMessage("[3/4] Обработано " + strconv.Itoa(count) + " из " + strconv.Itoa(len(data)))
			}

		}
	}
	wg.Wait()
	sentry.CaptureMessage("[3/4] Скрипт парсера картинок завершен!")
}

func main() {
	time.Sleep(time.Second * 10)
	fmt.Println("Images Started")
	for {
		scrapImages()
	}
}
