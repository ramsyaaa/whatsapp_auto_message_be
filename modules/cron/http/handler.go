package http

import (
	"encoding/json"
	"fmt"
	"go_whatsapp/helper"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type CronHandler struct {
}

func NewCronHandler() *CronHandler {
	return &CronHandler{}
}

func (h *CronHandler) FetchData(c *fiber.Ctx) error {
	baseURL := os.Getenv("BASE_URL")
	authHeader := os.Getenv("AUTH_HEADER")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/token?grant_type=password&username=%s&password=%s", baseURL, username, password), nil)
	req.Header.Set("Authorization", authHeader)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(respBody) == 0 {
		return c.Status(http.StatusOK).JSON(helper.APIResponse("Login Failed, Invalid Credentials", http.StatusBadRequest, "Error", nil))
	}
	var responseMap map[string]interface{}
	err = json.Unmarshal(respBody, &responseMap)
	if err != nil {
		return err
	}

	response := helper.APIResponse("Login Success", http.StatusOK, "OK", responseMap)
	return c.Status(http.StatusOK).JSON(response)
}
