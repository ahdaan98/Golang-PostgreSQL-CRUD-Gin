package controllers

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/ahdaan98/go-gorm-crud/config"
	"github.com/ahdaan98/go-gorm-crud/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Body struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
}

func CreateBook(c *gin.Context) {

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := CheckBookExist(config.DB, "string", "title", body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var resp models.Book
	query := `INSERT INTO books (title,author) VALUES(?,?) RETURNING *`

	if err := config.DB.Raw(query, body.Title, body.Author).Scan(&resp).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "failed to add book",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully Added the book",
		"data":    resp,
	})
}

func ListAllBook(c *gin.Context) {

	var books []models.Book
	query := `SELECT * FROM books ORDER BY id ASC`
	if err := config.DB.Raw(query).Scan(&books).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "failed to get all books",
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Successfully Retrieved All books",
		"books":   books,
	})
}

func UpdateBook(c *gin.Context) {
	bookID := c.Query("id")
	var body Body

	err := ValidateNumberParam(bookID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = CheckBookExist(config.DB, "id", bookID, Body{}); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	value, err := StrToInt(bookID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var count int
	query :=  `SELECT COUNT(*) FROM books WHERE id!=? AND title=?`
	if err:=config.DB.Raw(query,value,body.Title).Scan(&count).Error; err!=nil{
		c.JSON(400, gin.H{
			"error": "Failed to Check the book already exist",
		})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{
			"error": "Book with this Title Already exist",
		})
		return
	}


	var resp models.Book
	query = `UPDATE books SET title=?, author=? WHERE id=? RETURNING *`
	if err:=config.DB.Raw(query,body.Title,body.Author,value).Scan(&resp).Error; err!=nil{
		c.JSON(400, gin.H{
			"error": "Failed to Update book",
		})
		return
	}
	
	c.JSON(200, gin.H{
		"Message": "Successfully Updated the book",
		"updated_book":   resp,
	})
}

func DeleteBook(c *gin.Context) {
	bookID:=c.Query("id")

	err := ValidateNumberParam(bookID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = CheckBookExist(config.DB, "id", bookID, Body{}); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	value, err := StrToInt(bookID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := `DELETE FROM books WHERE id=?`

	if err := config.DB.Exec(query,value).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "failed to delete book",
		})
		return
	}

	c.JSON(400, gin.H{
		"message": "Successfully deleted the book",
	})
}

func GetBookByID(c *gin.Context){
	bookID:=c.Query("id")

	err := ValidateNumberParam(bookID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = CheckBookExist(config.DB, "id", bookID, Body{}); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	value, err := StrToInt(bookID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var book models.Book
	query := `SELECT * FROM books WHERE id=?`
	if err:=config.DB.Raw(query,value).Scan(&book).Error; err!=nil{
		c.JSON(400, gin.H{
			"error": "Failed to get book",
		})
		return
	}

	c.JSON(200,gin.H{
		"message":"Successfully got the book",
		"book":book,
	})
}

func ValidateNumberParam(val string) error {
	if val == ""{
		return errors.New("id cannot be nil")
	}

	v, err := StrToInt(val)

	if err != nil {
		return err
	}

	if reflect.TypeOf(v) != reflect.TypeOf(1) {
		return errors.New("value must be integer")
	}

	if v < 1 {
		return errors.New("values must be positibe number")
	}

	return nil
}

func CheckBookExist(DB *gorm.DB, t string, s string, b Body) error {
	if t == "string" {
		if s == "title" {
			query := "SELECT COUNT(*) FROM books WHERE title=?"
			var count int
			if err := DB.Raw(query, b.Title).Scan(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("book with this title already exist")
			}
		}
	}

	if t == "id" {
		id, err := StrToInt(s)
		if err != nil {
			return err
		}

		query := "SELECT COUNT(*) FROM books WHERE id=?"
		var count int
		if err := DB.Raw(query, id).Scan(&count).Error; err != nil {
			return err
		}
		if count < 1 {
			return errors.New("book with this id does not exist")
		}
	}

	return nil
}

func StrToInt(i string) (int, error) {
	if i == "" {
        return 0, errors.New("empty string cannot be converted to integer")
    }
	val, err := strconv.Atoi(i)
	if err != nil {
		return 0, errors.New("id can contain only number")
	}
	return val, nil
}
