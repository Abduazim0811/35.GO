package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func loadBooks() error {
	file, err := os.Open("books.json")
	defer file.Close()
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &books)
	return nil
}

func saveBooks() error {
	data, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile("books.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func main() {
	err := loadBooks()
	if err != nil {
		log.Fatal("Error loading books:", err)
	}

	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/book", getBook)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
