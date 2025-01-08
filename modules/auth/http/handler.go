package http

import (
	"bytes"
	"context"
	"go_whatsapp/helper"
	"image/png"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

type AuthHandler struct {
	client *whatsmeow.Client
}

func NewAuthHandler(client *whatsmeow.Client) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (h *AuthHandler) HandleQR(c *fiber.Ctx) error {
	if h.client.Store.ID == nil {
		qrChan, _ := h.client.GetQRChannel(context.Background())
		err := h.client.Connect()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to connect", http.StatusInternalServerError, "Error", nil))
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				// Generate QR code image
				qr, err := qrcode.New(evt.Code, qrcode.Medium)
				if err != nil {
					return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to generate QR code", http.StatusInternalServerError, "Error", nil))
				}

				// Convert to PNG
				var buf bytes.Buffer
				err = png.Encode(&buf, qr.Image(256))
				if err != nil {
					return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to encode QR code as PNG", http.StatusInternalServerError, "Error", nil))
				}

				c.Set("Content-Type", "image/png")
				return c.Send(buf.Bytes())
			}
		}
	}

	return c.Status(http.StatusOK).JSON(helper.APIResponse("Already logged in", http.StatusOK, "OK", nil))
}
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	if h.client.Store.ID == nil {
		return c.Status(http.StatusBadRequest).JSON(helper.APIResponse("Not logged in", http.StatusBadRequest, "Error", nil))
	}

	h.client.Disconnect()
	err := h.client.Store.Delete()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to logout", http.StatusInternalServerError, "Error", nil))
	}

	return c.Status(http.StatusOK).JSON(helper.APIResponse("Logged out successfully", http.StatusOK, "Success", nil))
}
