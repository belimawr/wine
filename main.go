package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Wine struct {
	Name        string
	Price       string
	Grape       string
	Description string
}

func wineGetter(out chan Wine, wg *sync.WaitGroup, url string) {
	doc, err := goquery.NewDocument("http://wine.com.br/" + url)

	if err != nil {
		panic(err)
	}

	name := doc.Find("#boxProduto > h1").Text()
	price := doc.Find("#boxProduto div.boxPreco > p").Text()
	grape := doc.Find("#paginaProduto > div.boxApresentacaoProduto > div.dadosAvancados > div > ul:nth-child(4) > li:nth-child(1) > span.valor").Text()
	description := doc.Find("#boxProduto div.comentarioSommelier > p").Text()

	out <- Wine{
		Name:        strings.Trim(name, "\n\t "),
		Price:       strings.Trim(price, "\n\t "),
		Grape:       strings.Trim(grape, "\n\t "),
		Description: strings.Trim(description, "\n\t "),
	}

	wg.Done()
}

func main() {
	doc, err := goquery.NewDocument("https://www.wine.com.br/browse.ep?cID=100851&filters=&pn=1&exibirEsgotados=false&listagem=horizontal&sorter=price-asc")

	if err != nil {
		panic(err.Error())
	}

	winesChan := make(chan Wine, 100)

	var wg sync.WaitGroup

	doc.Find("div.barraTitulo > h2 > a").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		if href, ok := s.Attr("href"); ok {
			go wineGetter(winesChan, &wg, href)
		}
	})

	wg.Wait()
	close(winesChan)

	for w := range winesChan {
		fmt.Println(w)
	}
}
