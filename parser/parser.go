package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/belimawr/wine/models"
)

// WinePage - Parser the details page of a wine
func WinePage(page *goquery.Selection) models.Wine {

	//w := models.Wine{}
	name := page.Find("#boxProduto > h1").Text()
	price := page.Find("#boxProduto div.boxPreco > p").Text()
	grape := page.Find("#paginaProduto > div.boxApresentacaoProduto > div.dadosAvancados > div > ul:nth-child(4) > li:nth-child(1) > span.valor").Text()
	description := page.Find("#boxProduto div.comentarioSommelier > p").Text()

	return models.Wine{
		Name:        strings.Trim(name, "\n\t "),
		Price:       strings.Trim(price, "\n\t "),
		Grape:       strings.Trim(grape, "\n\t "),
		Description: strings.Trim(description, "\n\t "),
	}
}

func ParseListing(page *goquery.Selection) []string {
	urls := []string{}

	page.Find("div.barraTitulo > h2 > a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			urls = append(urls, href)
		}
	})

	return urls
}
