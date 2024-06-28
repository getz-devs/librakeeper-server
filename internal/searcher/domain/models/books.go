package models

type BookSearchResult struct {
	ISBN      string
	Title     string
	Author    string
	Publisher string
	Year      string
}

type BooksSearchResult []*BookSearchResult
