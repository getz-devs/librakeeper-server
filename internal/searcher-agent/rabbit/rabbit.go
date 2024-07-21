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
	"time"
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

	books, err := h.scrapISBNFindBook(msg.GetIsbn())
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

func (h *Handler) scrapISBNFindBook(isbn string) ([]*bookModels.BookInShop, error) {
	const op = "rabbit.Handler.scrapISBNFindBook"
	log := h.log.With(slog.String("op", op), slog.String("isbn", isbn))
	preparedUrl := fmt.Sprintf(findBookUrlTemplate, isbn)
	var books []*bookModels.BookInShop

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		//colly.Async(true),
	)

	//Ignore the robot.txt
	c.IgnoreRobotsTxt = true
	// Time-out after 20 seconds.
	c.SetRequestTimeout(20 * time.Second)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	retryCount := 0
	maxRetryCount := 3

	c.OnResponse(func(r *colly.Response) {
		// print all headers
		if r.Headers.Get("Pragma") == "no-cache" && retryCount < maxRetryCount {
			time.Sleep(1 * time.Second)
			retryCount++
			r.Request.Visit(preparedUrl)
		}
	})

	c.OnHTML(
		"section.container.results",
		func(e *colly.HTMLElement) {
			e.ForEach("div.row.results__line", func(_ int, e *colly.HTMLElement) {
				book := &bookModels.BookInShop{}
				err := e.Unmarshal(book)
				if err != nil {
					return
				}
				if book.ImgUrl == "/images/camera.png" {
					book.ImgUrl = ""
				}

				books = append(books, book)
			})
		},
	)
	c.OnHTML("div.pagination__pages a:has(i.icon-angle-right)", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})
	c.OnRequest(func(r *colly.Request) {
		log.Info("visiting", slog.String("url", r.URL.String()))
	})

	err := c.Visit(preparedUrl)
	c.Wait()
	if err != nil {
		return nil, err
	}

	//for _, p := range books {
	//	fmt.Printf("%+v\n\n", p)
	//}
	if len(books) == 0 {
		return []*bookModels.BookInShop{}, nil
	}
	return books, nil
}
