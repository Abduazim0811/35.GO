package main

import (
	"encoding/json"
	"fmt"
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
	defer file.Close()

	err = json.NewDecoder(file).Decode(&books)
	if err != nil {
		return err
	}

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

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "Kitob topilmadi", http.StatusNotFound)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Malumotlarni o'qib bo'lmadi", http.StatusBadRequest)
		return
	}

	books = append(books, newBook)
	err = saveBooks()
	if err != nil {
		http.Error(w, "Kitob qo'shilmadi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Malumotlarni o'qib bo'lmadi", http.StatusBadRequest)
		return
	}

	for i, book := range books {
		if book.ID == id {
			books[i] = updatedBook
			err = saveBooks()
			if err != nil {
				http.Error(w, "Kitob yangilanmadi", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	http.Error(w, "Kitob topilmadi", http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			err = saveBooks()
			if err != nil {
				http.Error(w, "Kitob o'chirilmadi", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "Kitob topilmadi", http.StatusNotFound)
}

func main() {
	err := loadBooks()
	if err != nil {
		fmt.Println("Kitoblar yuklanmadi:", err)
		return
	}

	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/books/", getBook)
	http.HandleFunc("/books/add", addBook)
	http.HandleFunc("/books/update/", updateBook)
	http.HandleFunc("/books/delete/", deleteBook)

	fmt.Println("Server 8080 portda ishga tushirildi...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server ishlashda xatolik:", err)
	}
}
