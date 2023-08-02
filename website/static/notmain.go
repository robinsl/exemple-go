package static

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

type Book struct {
	ID    string
	Title string
	Page  int
}

type Books []Book

func NewBooks() Books {
	return Books{
		{
			ID:    "1",
			Title: "Book 1",
			Page:  100,
		},
		{
			ID:    "2",
			Title: "Book 2",
			Page:  200,
		},
		{
			ID:    "3",
			Title: "Book 3",
			Page:  300,
		},
	}
}

func (books Books) All() Books {
	return books
}

func (books Books) Get(id string) Book {
	for _, book := range books {
		if book.ID == id {
			return book
		}
	}

	return Book{}
}

func (books Books) Add(book Book) Books {
	books = append(books, book)
	return books
}

func (books Books) Delete(id string) Books {
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			return books
		}
	}

	return books
}

func (books Books) Update(id string, book Book) Books {
	for i, b := range books {
		if b.ID == id {
			books[i] = book
			return books
		}
	}

	return books
}

func main() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("No .env file found")
	}

	mongoConnectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/",
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"))

	_, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello World"))
	})

	books := NewBooks()
	router.Mount("/books", BookRouter(books))

	fs := http.FileServer(http.Dir("website/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", router)
}

func BookRouter(books Books) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(books.All())
	})

	router.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(books.Get(chi.URLParam(request, "id")))
	})

	router.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		var book Book
		json.NewDecoder(request.Body).Decode(&book)

		books = books.Add(book)

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(books)
	})

	router.Put("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		var book Book
		json.NewDecoder(request.Body).Decode(&book)

		books = books.Update(chi.URLParam(request, "id"), book)

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(books)
	})

	router.Delete("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		books = books.Delete(chi.URLParam(request, "id"))

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(books)
	})

	return router
}
