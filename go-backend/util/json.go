package util

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func JsonParsingError(c echo.Context) error {
	enrichedJson, err := json.Marshal(map[string]string{
		"message": "Error parsing request body. Please try again",
		"success": "false",
	})
	if err != nil {
		return err
	}
	return c.JSONBlob(400, enrichedJson)
}

func UserInputError(c echo.Context, message string) error {
	enrichedJson, err := json.Marshal(map[string]string{
		"message": message,
		"success": "false",
	})
	if err != nil {
		return err
	}
	return c.JSONBlob(400, enrichedJson)
}
