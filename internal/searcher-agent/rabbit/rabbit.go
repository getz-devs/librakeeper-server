package rabbit

import (
	"context"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/searcher-shared/domain/bookModels"
	rabbitDefines "github.com/getz-devs/librakeeper-server/lib/rabbit/getz.rabbitProto.v1"
	"github.com/gocolly/colly"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type Handler struct {
	log            *slog.Logger
	requestStorage RequestStorage
}

func New(log *slog.Logger, requestStorage RequestStorage) *Handler {
	return &Handler{
		log:            log,
		requestStorage: requestStorage,
	}
}

type RequestStorage interface {
	CompleteRequest(ctx context.Context, isbn string, books []*bookModels.BookInShop) error
	RejectRequest(ctx context.Context, isbn string) error
}

func (h *Handler) Handle(ctx context.Context, delivery amqp.Delivery) error {
	const op = "rabbit.Handler.Handle"
	log := h.log.With(slog.String("op", op))

	msg := &rabbitDefines.ISBNMessage{}
	if err := proto.Unmarshal(delivery.Body, msg); err != nil {
		log.Error("Error unmarshaling", err)
		return err
	}

	books, err := ScrapISBNFindBook(msg.GetIsbn())
	if err != nil {
		if err := h.requestStorage.RejectRequest(ctx, msg.GetIsbn()); err != nil {
			log.Error("Error rejecting request", err)
		}
		return err
	}

	if err := h.requestStorage.CompleteRequest(ctx, msg.GetIsbn(), books); err != nil {
		log.Error("Error completing request", err)
		return err
	}

	return nil
}

const findBookUrlTemplate = "https://www.findbook.ru/search/d1?isbn=%s&r=0&s=1&viewsize=15&startidx=0"

func ScrapISBNFindBook(isbn string) ([]*bookModels.BookInShop, error) {
	var books []*bookModels.BookInShop

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

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
		return nil, err
	}

	for _, p := range books {
		fmt.Printf("%+v\n\n", p)
	}
	return books, nil
}
