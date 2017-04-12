package main

import (
	"database/sql"
	"log"
	"sync"

	"fmt"

	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/belimawr/wine/models"
	"github.com/belimawr/wine/parser"
	"github.com/belimawr/wine/store"
	_ "github.com/mattn/go-sqlite3"
)

func wineGetter(i int, in <-chan string, out chan<- models.Wine) {
	for {
		url, more := <-in

		if !more {
			break
		}

		doc, err := goquery.NewDocument("http://wine.com.br/" + url)

		if err != nil {
			log.Printf("Error reading Wine page: %q", err)
		}

		out <- parser.WinePage(doc.Selection)

		fmt.Fprintf(os.Stderr, " %02d ", i)
	}
	log.Print("Worker Done!")

}

func main() {
	db, err := sql.Open("sqlite3", "./sqlite3.db")

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1)

	createWineTable := `CREATE TABLE wine(id INTEGER NOT NULL PRIMARY KEY,
                                         name TEXT,
                                         price REAL,
                                         deal REAL,
                                         grape TEXT,
                                         description TEXT,
                                         pairing TEXT,
                                         crawled_at TIMESTAMP,
                                         error TEXT);`

	_, err = db.Exec(createWineTable)

	if err != nil {
		log.Print(err)
	}

	store := store.NewSQLiteStore(db)

	urlPattern := "https://www.wine.com.br/browse.ep?cID=100851&filters=&pn=%d&exibirEsgotados=false&listagem=horizontal&sorter=price-asc"

	winesChan := make(chan models.Wine, 100)
	urlChan := make(chan string, 100)

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		go wineGetter(i, urlChan, winesChan)

		go func(winesChan <-chan models.Wine, wg *sync.WaitGroup) {
			for {
				w, more := <-winesChan
				if !more {
					break
				}
				store.PutWine(w)
				wg.Done()
			}
			log.Println("Consumer done")
		}(winesChan, &wg)
	}

	i := 0
	wgc := 0
	for {
		i++
		doc, err := goquery.NewDocument(fmt.Sprintf(urlPattern, i))

		if err != nil {
			fmt.Print(err.Error())
		}
		log.Printf("Got page: %d", i)

		urls := parser.ParseListing(doc.Find("body"))

		fmt.Printf("\nPage: %d, urls: %d\n\n", i, len(urls))

		if len(urls) == 0 {
			break
		}

		for _, url := range urls {
			wg.Add(1)
			urlChan <- url
			wgc++
			fmt.Print(".")
		}
	}

	fmt.Printf("\n\nwgc: %d\n\n", wgc)
	wg.Wait()

	close(urlChan)
	close(winesChan)

}
