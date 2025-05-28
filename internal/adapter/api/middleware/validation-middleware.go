package middleware

import (
	"codetest/internal/adapter/api/presenter"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type BindingType int

const (
	BindJSON BindingType = iota
	BindQuery
	BindForm
	BindUri
	BindHeader
	BindMultipartForm
)

func ValidationMiddleware(T any, bindingType BindingType) gin.HandlerFunc {
	return func(c *gin.Context) {
		tType := reflect.TypeOf(T)

		// If T is not a pointer, create a pointer to it
		if tType.Kind() != reflect.Ptr {
			tType = reflect.PointerTo(tType)
		}

		// Create a new instance of the type
		val := reflect.New(tType.Elem()).Interface()

		if err := getBindingError(c, val, bindingType); err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := make(map[string]string)
				for _, e := range errors {
					field, _ := tType.Elem().FieldByName(e.Field())
					jsonTag := field.Tag.Get(getTagName(bindingType))
					errorMessages[jsonTag] = formatValidationMessages(jsonTag, e.Tag(), e.Param())
				}

				c.JSON(400, presenter.JsonResponseWithoutPagination{
					Success: false,
					Data:    nil,
					Error:   errorMessages,
				})
				c.Abort()
				return
			}

			c.JSON(400, presenter.JsonResponseWithoutPagination{
				Success: false,
				Data:    nil,
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("validatedRequest", val)
		c.Next()
	}
}

func getBindingError(c *gin.Context, val any, bindingType BindingType) error {
	var err error

	switch bindingType {
	case BindJSON:
		err = c.ShouldBindJSON(val)
	case BindQuery:
		err = c.ShouldBindQuery(val)
	case BindForm:
		err = c.ShouldBind(val)
	case BindUri:
		err = c.ShouldBindUri(val)
	case BindHeader:
		err = c.ShouldBindHeader(val)
	case BindMultipartForm:
		err = c.ShouldBindWith(val, binding.FormMultipart)
	default:
		err = c.ShouldBindJSON(val)
	}

	return err
}

func getTagName(bindingType BindingType) string {
	var tagName string

	switch bindingType {
	case BindJSON:
		tagName = "json"
	case BindQuery, BindForm, BindMultipartForm:
		tagName = "form"
	case BindUri:
		tagName = "uri"
	case BindHeader:
		tagName = "header"
	default:
		tagName = "json"
	}

	return tagName
}

func formatValidationMessages(field, tag, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + param + " characters long"
	case "max":
		return field + " must be at most " + param + " characters long"
	case "gte":
		return field + " must be greater than or equal to " + param
	case "lte":
		return field + " must be less than or equal to " + param
	case "eq":
		return field + " must be equal to " + param
	case "ne":
		return field + " must not be equal to " + param
	case "gt":
		return field + " must be greater than " + param
	case "lt":
		return field + " must be less than " + param
	case "len":
		return field + " must be exactly " + param + " characters long"
	case "oneof":
		return field + " must be one of the following values: " + param
	case "notoneof":
		return field + " must not be one of the following values: " + param
	case "url":
		return field + " must be a valid URL"
	case "uuid":
		return field + " must be a valid UUID"
	case "json":
		return field + " must be a valid JSON"
	default:
		return field + " is invalid"
	}
}
