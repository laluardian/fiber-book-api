package controllers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/laluardian/fiber-book-api/config"
	"github.com/laluardian/fiber-book-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FetchBooks(c *fiber.Ctx) error {
	bookCollection := config.MI.DB.Collection(os.Getenv("BOOK_COLLECTION"))

	query := bson.D{{}}

	cursor, err := bookCollection.Find(c.Context(), query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	var books []models.Book = make([]models.Book, 0)

	err = cursor.All(c.Context(), &books)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"books": books,
		},
	})
}

func InsertBook(c *fiber.Ctx) error {
	bookCollection := config.MI.DB.Collection(os.Getenv("BOOK_COLLECTION"))

	data := new(models.Book)

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	data.ID = nil
	f := false
	data.IsReading = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	result, err := bookCollection.InsertOne(c.Context(), data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert book",
			"error":   err,
		})
	}

	book := &models.Book{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	bookCollection.FindOne(c.Context(), query).Decode(book)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"book": book,
		},
	})
}

func FetchBook(c *fiber.Ctx) error {
	bookCollection := config.MI.DB.Collection(os.Getenv("BOOK_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse Id",
			"error":   err,
		})
	}

	book := &models.Book{}

	query := bson.D{{Key: "_id", Value: id}}

	err = bookCollection.FindOne(c.Context(), query).Decode(book)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Book not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"book": book,
		},
	})
}

func UpdateBook(c *fiber.Ctx) error {
	bookCollection := config.MI.DB.Collection(os.Getenv("BOOK_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	data := new(models.Book)
	err = c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	var dataToUpdate bson.D

	if data.Title != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "title", Value: data.Title})
	}

	if data.IsReading != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "completed", Value: data.IsReading})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	update := bson.D{
		{Key: "$set", Value: dataToUpdate},
	}

	err = bookCollection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Book Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update book",
			"error":   err,
		})
	}

	book := &models.Book{}

	bookCollection.FindOne(c.Context(), query).Decode(book)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"book": book,
		},
	})
}

func DeleteBook(c *fiber.Ctx) error {
	bookCollection := config.MI.DB.Collection(os.Getenv("BOOK_COLLECTION"))

	paramID := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	err = bookCollection.FindOneAndDelete(c.Context(), query).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Book Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete book",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
