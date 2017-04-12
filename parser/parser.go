package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/belimawr/wine/models"
)

var whiteSpaceReplacer = strings.NewReplacer("\n", "", "\t", "")
var priceReplacer = strings.NewReplacer(".", "", " ", "", "\n", "", "\t", "", "R$", "", ",", ".")

// WinePage - Parses the page of a wine and returns a Wine
func WinePage(page *goquery.Selection) models.Wine {

	name := page.Find("#boxProduto > h1").Text()
	price := page.Find("#boxProduto div.boxPreco > p").Text()
	grape := page.Find("#paginaProduto > div.boxApresentacaoProduto > div.dadosAvancados > div > ul:nth-child(4) > li:nth-child(1) > span.valor").Text()
	pairing := page.Find("#paginaProduto > div.boxApresentacaoProduto > div.dadosAvancados > div > ul:nth-child(8) > li:nth-child(4) > span.valor").Text()
	description := page.Find("#boxProduto div.comentarioSommelier > p").Text()
	deal := "0.0"

	name = whiteSpaceReplacer.Replace(name)
	price = priceReplacer.Replace(price)
	grape = whiteSpaceReplacer.Replace(grape)
	description = whiteSpaceReplacer.Replace(description)
	pairing = whiteSpaceReplacer.Replace(pairing)

	wine := models.Wine{
		Name:        name,
		Description: description,
		Grape:       grape,
		Pairing:     pairing,
	}

	if strings.HasPrefix(price, "De") {
		price = strings.Replace(price, "De", "", -1)

		split := strings.Split(price, "por")

		if len(split) == 2 {
			price = split[0]
			deal = split[1]
		}
	}

	if price != "ProdutoIndisponível" {
		if floatPrice, err := strconv.ParseFloat(price, 64); err != nil {
			log.Printf("Could not convert %q to float. Error: %s", price, err.Error())
			wine.Error = fmt.Sprintf("%q, %q", wine.Error, err.Error())
		} else {
			wine.Price = floatPrice
		}

		if floatDeal, err := strconv.ParseFloat(deal, 64); err != nil {
			log.Printf("Could not convert %q to float. Error: %s", deal, err.Error())
			wine.Error = fmt.Sprintf("%q, %q", wine.Error, err.Error())
		} else {
			wine.Deal = floatDeal
		}
	}

	return wine
}

// ParseListing - Parses listing pages and returns a list of liks to be visited
func ParseListing(page *goquery.Selection) []string {
	urls := []string{}

	page.Find("div.barraTitulo > h2 > a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			urls = append(urls, href)
		}
	})

	return urls
}
