package searcher_agent_cmd

import (
	"fmt"
	colly "github.com/gocolly/colly"
)

func main() {
	// TODO
	err := scrapISBNFindBook("9785206000344")
	if err != nil {
		return
	}
}

var (
	const findbookUrlTemplate = "https://www.findbook.ru/search/d1?isbn=%s"
)
func scrapISBNFindBook(isbn string) error {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	c.OnHTML(
		"section.container section.results",
		func(e *colly.HTMLElement) {
			//fmt.Println(e.Attr("href"))
			results_book := e.ChildText("")

			// TODO
		},
	)

	preparedUrl := fmt.Sprintf(findbookUrlTemplate, isbn)
	err := c.Visit(preparedUrl)
	if err != nil {
		return err
	}


	return ""
}

type book_result struct {
	// TODO
}
