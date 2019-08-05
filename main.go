package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Book struct (MODEL)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//Author Struct
type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

//init books var as a slice Book struct
var books []Book

//Get all books
func getBooks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(books)
}

//Get single book
func getBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // get params
	//loop through books and find the one that matches id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(res).Encode(item)
			return
		}
	}
	json.NewEncoder(res).Encode(&Book{})
}

//Create single book
func createBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) // mock ID - not safe
	books = append(books, book)
	json.NewEncoder(res).Encode(book)
}

//Update single book
func updateBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			var book Book
			_ = json.NewDecoder(req.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(res).Encode(book)
			return
		}
	}
	json.NewEncoder(res).Encode(books)
}

//Delete single book
func deleteBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(res).Encode(books)
}

func main() {
	//init router
	router := mux.NewRouter()

	//Mock data - @todo - implement DB
	books = append(books, Book{ID: "0", Isbn: "1452543", Title: "Book One", Author: &Author{FirstName: "John", LastName: "Doe"}})
	books = append(books, Book{ID: "1", Isbn: "68563", Title: "Book Two", Author: &Author{FirstName: "Steve", LastName: "Smith"}})
	books = append(books, Book{ID: "2", Isbn: "53464567", Title: "Book Three", Author: &Author{FirstName: "George", LastName: "Martin"}})

	//route handlers / endpoints
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
