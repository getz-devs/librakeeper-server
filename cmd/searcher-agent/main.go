package main

import (
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/searcher-agent/app"
	"github.com/getz-devs/librakeeper-server/internal/searcher-agent/config"
	"github.com/getz-devs/librakeeper-server/lib/prettylog"
	colly "github.com/gocolly/colly"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//err := scrapISBNFindBook("9785206000344")
	//if err != nil {
	//	panic(err)
	//}

	//const rabbitUrl = "amqp://guest:guest@192.168.1.161:5672/"

	cfg := config.MustLoad()

	log := prettylog.SetupLogger(cfg.Env)

	log.Info("starting ...",
		slog.String("env", cfg.Env),
		slog.Any("config", cfg),
	)

	application := app.New(cfg.ConnectUrl, cfg.QueueName, log)
	go application.AppRabbit.MustRun()

	// --------------------------- Register stop signal ---------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// --------------------------- Wait for stop signal ---------------------------
	sign := <-stop

	log.Info("shutting down ...",
		slog.String("signal", sign.String()),
	)

	application.AppRabbit.Close()

	//application.

	log.Info("application fully stopped")
}

const findBookUrlTemplate = "https://www.findbook.ru/search/d1?isbn=%s"

func scrapISBNFindBook(isbn string) error {
	var books []*bookInShop

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	c.OnHTML(
		"section.container.results",
		func(e *colly.HTMLElement) {
			//fmt.Println(e.Attr("href"))
			//fmt.Println("Founded \n\n", e)
			e.ForEach("div.row.results__line", func(_ int, e *colly.HTMLElement) {
				book := &bookInShop{}
				e.Unmarshal(book)

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

type bookInShop struct {
	Title      string `selector:"div.results__book-name > a"`
	Author     string `selector:"div.results__authors"`
	Publishing string `selector:"div.results__publishing"`
	ImgUrl     string `selector:"a.results__image > img" attr:"src"`
	ShopName   string `selector:"div.results__shop-name > a"`
}
