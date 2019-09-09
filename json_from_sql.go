package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Book Model
type Book struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	PublicationDate string  `json:"publicationDate"`
	Author          *Author `json:"author"`
}

//Author Model
type Author struct {
	ID        string `json:"authorID"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applocation/json")
	params := mux.Vars(r)

	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	db, err := sql.Open("mysql", "mohammad:root@tcp(127.0.0.1:3306)/test")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT Books.ID, Books.name, Books.pubDate, Authors.ID, Authors.firstname, Authors.lastname FROM Books, Authors, AuthorsHaveBooks WHERE (Books.ID = AuthorsHaveBooks.book_ID AND Authors.ID = AuthorsHaveBooks.author_ID)")

	if err != nil {
		fmt.Println("Error running Query!")
		fmt.Println(err)
		return
	}

	defer rows.Close()

	i := 0

	for rows.Next() {
		var bookID string
		var bookName string
		var bookPublicationDate string
		var bookAuthorID string
		var bookAuthorFirstname string
		var bookAuthorLastname string

		rows.Scan(&bookID, &bookName, &bookPublicationDate, &bookAuthorID, &bookAuthorFirstname, &bookAuthorLastname)
		author := Author{ID: bookAuthorID, Firstname: bookAuthorFirstname, Lastname: bookAuthorLastname}
		books = append(books,
			Book{ID: bookID, Name: bookName, PublicationDate: bookPublicationDate, Author: &author})

		i++
	}

	//Initialize Router
	router := mux.NewRouter()
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", addBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe(":8000", router)

}
