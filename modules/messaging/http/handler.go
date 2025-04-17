package http

import (
	"context"
	"fmt"
	"go_whatsapp/helper"
	"go_whatsapp/modules/messaging/service"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type MessagingHandler struct {
	client  *whatsmeow.Client
	service service.MessageService
}

func NewMessagingHandler(client *whatsmeow.Client, service service.MessageService) *MessagingHandler {
	return &MessagingHandler{
		client:  client,
		service: service,
	}
}

type LinkPreview struct {
	MatchedText   string
	CanonicalURL  string
	Title         string
	Description   string
	JPEGThumbnail []byte
}

type SendMessagePayload struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func (h *MessagingHandler) HandleSendMessage(c *fiber.Ctx) error {
	// Check if client is logged in
	if h.client.Store.ID == nil {
		return c.Status(400).JSON(helper.APIResponse("Please login first", http.StatusBadRequest, "ERROR", nil))
	}

	// Parse the payload
	payload := new(SendMessagePayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(helper.APIResponse("Invalid payload", http.StatusBadRequest, "ERROR", nil))
	}

	// Format number (should start with '62' for Indonesia)
	number := payload.Number
	if !strings.HasPrefix(number, "62") {
		if strings.HasPrefix(number, "0") {
			number = "62" + number[1:] // Replace '0' with '62' for Indonesia
		} else if strings.HasPrefix(number, "+62") {
			number = number[1:] // Remove leading '+' for international numbers
		}
	}

	// Create the JID for the recipient
	recipient := types.JID{
		User:   number,
		Server: "s.whatsapp.net", // Server for WhatsApp
	}

	// Check if message contains URL and get preview if it does
	var linkPreview *waProto.Message
	if containsURL(payload.Message) {
		fmt.Println("URL FOUND")
		// Fetch the link preview
		preview, err := fetchLinkPreview(payload.Message)
		if err == nil && preview != nil {
			linkPreview = &waProto.Message{
				ExtendedTextMessage: &waProto.ExtendedTextMessage{
					Text:          proto.String(payload.Message),
					MatchedText:   proto.String(preview.MatchedText),
					CanonicalURL:  proto.String(preview.CanonicalURL),
					Title:         proto.String(preview.Title),
					Description:   proto.String(preview.Description),
					JPEGThumbnail: preview.JPEGThumbnail,
				},
			}
		}
	}

	// Create the message with preview if available
	var msg *waProto.Message
	if linkPreview != nil {
		msg = linkPreview
	} else {
		msg = &waProto.Message{
			Conversation: proto.String(payload.Message),
		}
	}

	// Set a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ensure client is connected
	if !h.client.IsConnected() {
		err := h.client.Connect()
		if err != nil {
			return c.Status(500).JSON(helper.APIResponse(fmt.Sprintf("Failed to connect to WhatsApp: %v", err), http.StatusInternalServerError, "ERROR", nil))
		}
	}

	// Send the message
	_, err := h.client.SendMessage(ctx, recipient, msg)
	if err != nil {
		return c.Status(500).JSON(helper.APIResponse(fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError, "ERROR", nil))
	}

	// Save message to database
	message, err := h.service.SaveMessage(ctx, number, payload.Message)
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to save message to database: %v\n", err)
	}

	return c.Status(200).JSON(helper.APIResponse("Message sent successfully", http.StatusOK, "SUCCESS", fiber.Map{
		"message_id": message.ID,
		"to":         number,
		"sent_at":    message.SentAt,
	}))
}

// Get recent messages
func (h *MessagingHandler) GetRecentMessages(c *fiber.Ctx) error {
	ctx := context.Background()
	
	// Get recent messages from database
	messages, err := h.service.GetRecentMessages(ctx, 10) // Get last 10 messages
	if err != nil {
		return c.Status(500).JSON(helper.APIResponse(fmt.Sprintf("Failed to get recent messages: %v", err), http.StatusInternalServerError, "ERROR", nil))
	}
	
	return c.Status(200).JSON(helper.APIResponse("Recent messages retrieved successfully", http.StatusOK, "SUCCESS", messages))
}

// Utility function to check for URLs in the message
func containsURL(message string) bool {
	return strings.Contains(message, "http://") || strings.Contains(message, "https://")
}

// Function to fetch link preview using goquery for scraping
func fetchLinkPreview(message string) (*LinkPreview, error) {
	// Extract URL from message using simple regex
	urlStart := strings.Index(message, "http")
	if urlStart == -1 {
		return nil, fmt.Errorf("no URL found in message")
	}

	// Find end of URL (space, newlines, or end of string)
	urlEnd := -1
	for i, char := range message[urlStart:] {
		if char == ' ' || char == '\n' {
			urlEnd = i
			break
		}
	}

	var url string
	if urlEnd == -1 {
		url = message[urlStart:]
	} else {
		url = message[urlStart : urlStart+urlEnd]
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract metadata
	title := doc.Find("meta[property='og:title']").First().Text()
	description := doc.Find("meta[property='og:description']").AttrOr("content", "")
	fmt.Println(title, description)

	// Get favicon/thumbnail
	var thumbnail []byte
	thumbnailURL := doc.Find("meta[property='og:image']").AttrOr("content", "")
	if thumbnailURL != "" {
		if !strings.HasPrefix(thumbnailURL, "http") {
			thumbnailURL = url + thumbnailURL
		}
		resp, err := client.Get(thumbnailURL)
		if err == nil {
			thumbnail, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}

	preview := &LinkPreview{
		MatchedText:   url,
		CanonicalURL:  url,
		Title:         title,
		Description:   description,
		JPEGThumbnail: thumbnail,
	}

	return preview, nil
}
