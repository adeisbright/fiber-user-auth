package middleware

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateCreateUser(c *fiber.Ctx) error {
	type User struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=1"`
		Password string `json:"password" validate:"required,min=8,max=60"`
	}

	// type IError struct {
	// 	Field string
	// 	Tag   string
	// 	Value string
	// }

	// var errors []*IError

	body := User{}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable to parse user creation data")
	}

	Validator := validator.New()

	err = Validator.Struct(body)
	if err != nil {
		errMsgs := make([]string, 0)

		for _, err := range err.(validator.ValidationErrors) {
			// var elem IError
			// elem.Field = err.Field()
			// elem.Tag = err.Kind().String() // Export struct tag
			// elem.Value = err.Param()       // Export field value

			// errors = append(errors, &elem)
			errMsgs = append(errMsgs, fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", err.Field(), err.Param(), err.Tag()))
		}
		fmt.Println(err.Error())

		//return c.Status(fiber.StatusBadRequest).JSON(errors)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}
	return c.Next()
}
