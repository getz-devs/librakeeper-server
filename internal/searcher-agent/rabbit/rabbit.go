package rabbit

import (
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/searcher-shared/domain/bookModels"
	rabbitDefines "github.com/getz-devs/librakeeper-server/lib/rabbit/getz.rabbitProto.v1"
	"github.com/gocolly/colly"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

//type Parser struct {
//	log *slog.Logger
//}
//
//func New(log *slog.Logger) *Parser {
//	const op
//	return &Parser{
//		log: log,
//	}
//}

func Handler(delivery amqp.Delivery, log *slog.Logger) error {
	const op = "rabbit.Handler"
	log = log.With(slog.String("op", op))

	msg := &rabbitDefines.ISBNMessage{}
	err := proto.Unmarshal(delivery.Body, msg)
	if err != nil {
		return err
	}

	err = scrapISBNFindBook(msg.GetIsbn())
	if err != nil {
		return err
	}

	return nil
}

const findBookUrlTemplate = "https://www.findbook.ru/search/d1?isbn=%s"

func scrapISBNFindBook(isbn string) error {
	var books []*bookModels.BookInShop

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	c.OnHTML(
		"section.container.results",
		func(e *colly.HTMLElement) {

			e.ForEach("div.row.results__line", func(_ int, e *colly.HTMLElement) {
				book := &bookModels.BookInShop{}
				err := e.Unmarshal(book)
				if err != nil {
					return
				}

				books = append(books, book)
			})

		},
	)

	preparedUrl := fmt.Sprintf(findBookUrlTemplate, isbn)
	err := c.Visit(preparedUrl)
	if err != nil {
		return err
	}

	for _, p := range books {
		fmt.Printf("%+v\n\n", p)
	}
	return nil
}
