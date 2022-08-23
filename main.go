package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

type book struct{
	ID 		string `json:"id"` 
	Title 	string `json:"title"`
	Author 	string `json:"author"`
	Quantity int64 `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context){
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context){
	id := c.Param("id")
	books, err := getBooksById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound , gin.H{
			"message": "book not found!",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

func checkOutBook(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"Misssing id querry param"})
		return
	}

	book, err := getBooksById(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound , gin.H{
			"message": "book not found!",
		})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest , gin.H{
			"message": "book not availavle!",
		})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func unCheckBook(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"Misssing id querry param"})
		return
	}

	book, err := getBooksById(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound , gin.H{
			"message": "book not found!",
		})
		return
	}

	if book.Quantity < 0 {
		c.IndentedJSON(http.StatusBadRequest , gin.H{
			"message": "book not availavle!",
		})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBooksById(id string)(*book, error){
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil,errors.New("book not found")
}

func createBook(c *gin.Context){
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default() 
	router.GET("/book", getBooks)
	router.GET("/book/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkOutBook)
	router.PATCH("/uncheckout", unCheckBook)
	router.Run("localhost:7000")
}