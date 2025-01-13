package http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_whatsapp/helper"
	"go_whatsapp/modules/broadcast/service"

	"github.com/PuerkitoBio/goquery"

	"github.com/gofiber/fiber/v2"
	"github.com/tealeg/xlsx"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"google.golang.org/protobuf/proto"
)

type BroadcastHandler struct {
	service service.BroadcastService
	client  *whatsmeow.Client
}

type LinkPreview struct {
	MatchedText   string
	CanonicalURL  string
	Title         string
	Description   string
	JPEGThumbnail []byte
}

type SendMessagePayload struct {
	BroadcastID int `json:"broadcast_id"`
}

func NewBroadcastHandler(service service.BroadcastService, client *whatsmeow.Client) *BroadcastHandler {
	return &BroadcastHandler{service: service, client: client}
}

func (h *BroadcastHandler) CreateBroadcast(c *fiber.Ctx) error {
	ctx := context.Background()
	type Request struct {
		ClientInfo       string `json:"client_info"`
		BroadcastPlanAt  string `json:"broadcast_plan_at"`
		BroadcastMessage string `json:"broadcast_message"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		response := helper.APIResponse("Invalid request data", http.StatusBadRequest, "ERROR", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Menghasilkan kode broadcast menggunakan fungsi generateBroadcastCode
	broadcastCode := generateBroadcastCode()

	broadcastPlanAt, err := time.Parse(time.RFC3339, req.BroadcastPlanAt)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(helper.APIResponse("Invalid broadcast plan at format", http.StatusBadRequest, "ERROR", nil))
	}
	data := map[string]interface{}{
		"client_info":       req.ClientInfo,
		"broadcast_plan_at": broadcastPlanAt.Format(time.RFC3339),
		"broadcast_message": req.BroadcastMessage,
		"created_at":        time.Now().Format(time.RFC3339),
		"broadcast_code":    broadcastCode,
	}

	broadcast, err := h.service.CreateBroadcast(ctx, data)
	if err != nil {
		response := helper.APIResponse("Failed to create broadcast", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.APIResponse("Broadcast created successfully", http.StatusOK, "OK", broadcast)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *BroadcastHandler) ImportRecipient(c *fiber.Ctx) error {
	ctx := context.Background()
	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIResponse("Gagal mengimpor file", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		response := helper.APIResponse("Gagal membuka file", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	defer src.Close()

	// Read the file into a slice of bytes
	fileData, err := ioutil.ReadAll(src)
	if err != nil {
		response := helper.APIResponse("Gagal membaca file", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Parse the file
	xlFile, err := xlsx.OpenBinary(fileData)
	if err != nil {
		response := helper.APIResponse("Gagal membuka file excel", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Iterate over the rows
	var recipients []map[string]interface{}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 { // Skip the first row because it is header
				continue
			}
			// Validate the number and change it if necessary
			cell := row.Cells[0]
			if cell != nil {
				number := cell.String()
				if strings.HasPrefix(number, "08") || strings.HasPrefix(number, "+62") {
					number = "62" + number[1:]
				}
				// Add the recipient to the data
				recipients = append(recipients, map[string]interface{}{"whatsapp_number": number})
			}
		}
	}

	// Mendapatkan broadcast_id dari form input
	broadcastId, err := strconv.Atoi(c.FormValue("broadcast_id"))
	if err != nil {
		response := helper.APIResponse("Invalid broadcast id", http.StatusBadRequest, "ERROR", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Import the recipients
	_, err = h.service.ImportRecipient(ctx, broadcastId, recipients)
	if err != nil {
		response := helper.APIResponse("Gagal mengimpor penerima", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.APIResponse("Penerima berhasil diimpor", http.StatusOK, "OK", nil)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *BroadcastHandler) ImportPecatuRecipient(c *fiber.Ctx) error {
	ctx := context.Background()
	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIResponse("Gagal mengimpor file", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		response := helper.APIResponse("Gagal membuka file", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	defer src.Close()

	// Read the file into a slice of bytes
	fileData, err := ioutil.ReadAll(src)
	if err != nil {
		response := helper.APIResponse("Gagal membaca file", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Parse the file
	xlFile, err := xlsx.OpenBinary(fileData)
	if err != nil {
		response := helper.APIResponse("Gagal membuka file excel", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	// Iterate over the rows
	var recipients []map[string]interface{}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 { // Skip the first row because it is header
				continue
			}
			// Validate the number and change it if necessary
			cellNumber := row.Cells[0]
			cellName := row.Cells[1]
			cellIdentifier := row.Cells[2]
			if cellNumber != nil && cellName != nil && cellIdentifier != nil {
				number := cellNumber.String()
				name := cellName.String()
				identifier := cellIdentifier.String()
				if strings.HasPrefix(number, "08") || strings.HasPrefix(number, "+62") {
					number = "62" + number[1:]
				}
				// Add the recipient to the data
				recipients = append(recipients, map[string]interface{}{"whatsapp_number": number, "name": name, "identifier": identifier})
			}
		}
	}
	// Mendapatkan broadcast_id dari form input
	broadcastId, err := strconv.Atoi(c.FormValue("broadcast_id"))
	if err != nil {
		response := helper.APIResponse("Invalid broadcast id", http.StatusBadRequest, "ERROR", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// Import the recipients
	_, err = h.service.ImportPecatuRecipient(ctx, broadcastId, recipients)
	if err != nil {
		response := helper.APIResponse("Gagal mengimpor penerima", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.APIResponse("Penerima berhasil diimpor", http.StatusOK, "OK", nil)
	return c.Status(http.StatusOK).JSON(response)
}

func generateBroadcastCode() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	broadcastCode := make([]byte, 6)
	for i := range broadcastCode {
		broadcastCode[i] = chars[rand.Intn(len(chars))]
	}
	return string(broadcastCode)
}

func (h *BroadcastHandler) BroadcastDetail(c *fiber.Ctx) error {
	ctx := context.Background()
	broadcastID, err := strconv.Atoi(c.Params("broadcast_id"))
	if err != nil {
		response := helper.APIResponse("Invalid broadcast id", http.StatusBadRequest, "ERROR", nil)
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	detail, err := h.service.BroadcastDetail(ctx, broadcastID)
	if err != nil {
		response := helper.APIResponse("Failed to get broadcast detail", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := helper.APIResponse("Broadcast detail retrieved successfully", http.StatusOK, "OK", detail)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *BroadcastHandler) HandleSendBroadcast(c *fiber.Ctx) error {
	// Check if client is logged in
	if h.client.Store.ID == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Please login first",
		})
	}

	// Parse the payload
	payload := new(SendMessagePayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	// Get the broadcast ID from the payload
	broadcastID := payload.BroadcastID

	// Check if there are any recipients in the broadcast
	ctx := context.Background()
	hasRecipients, err := h.service.IsAnyRecipientInBroadcast(ctx, broadcastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to check recipients",
		})
	}
	if !hasRecipients {
		return c.Status(200).JSON(fiber.Map{
			"message": "No recipients to send broadcast",
		})
	}

	// Update broadcast status to "Starting"
	_, err = h.service.UpdateBroadcastStatus(ctx, broadcastID, "Starting")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update broadcast status",
		})
	}

	// Get the broadcast message
	broadcastMessage, err := h.service.GetBroadcastMessage(ctx, broadcastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get broadcast message",
		})
	}

	// Get recipients by broadcast ID
	recipients, err := h.service.GetAllRecipientByBroadcastID(ctx, broadcastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get recipients",
		})
	}

	// Send an immediate response
	response := fiber.Map{
		"message": "Message sending process has started",
	}
	if err := c.JSON(response); err != nil {
		return err // Handle any error in sending the response
	}

	// Start the background process
	go func() {
		ch := make(chan error, len(recipients["recipients"].([]map[string]interface{})))

		for _, recipient := range recipients["recipients"].([]map[string]interface{}) {
			recipientID, ok := recipient["id"].(int)
			if !ok {
				ch <- fmt.Errorf("recipient ID is not an integer")
				return
			}
			whatsappNumber, ok := recipient["whatsapp_number"].(string)
			if !ok {
				ch <- fmt.Errorf("WhatsApp number is not a string")
				return
			}

			// Ensure client is connected
			if !h.client.IsConnected() {
				err := h.client.Connect()
				if err != nil {
					ch <- fmt.Errorf("failed to connect to WhatsApp: %v", err)
					return
				}
			}

			// Prepare recipient JID and message
			recipientJID := types.JID{
				User:   whatsappNumber,
				Server: "s.whatsapp.net",
			}
			var msg *waProto.Message
			if containsURL(broadcastMessage["broadcast_message"].(string)) {
				preview, err := fetchLinkPreview(broadcastMessage["broadcast_message"].(string))
				if err == nil && preview != nil {
					msg = &waProto.Message{
						ExtendedTextMessage: &waProto.ExtendedTextMessage{
							Text:          proto.String(broadcastMessage["broadcast_message"].(string)),
							MatchedText:   proto.String(preview.MatchedText),
							CanonicalURL:  proto.String(preview.CanonicalURL),
							Title:         proto.String(preview.Title),
							Description:   proto.String(preview.Description),
							JPEGThumbnail: preview.JPEGThumbnail,
						},
					}
				}
			}
			if msg == nil {
				msg = &waProto.Message{
					Conversation: proto.String(broadcastMessage["broadcast_message"].(string)),
				}
			}

			// Send message
			_, err = h.client.SendMessage(ctx, recipientJID, msg)
			if err != nil {
				// Jika pengiriman pesan gagal, lanjutkan ke penerima berikutnya dan update status menjadi gagal
				_, errUpdate := h.service.UpdateRecipientBroadcastStatus(ctx, recipientID, broadcastID, "Failed")
				if errUpdate != nil {
					ch <- fmt.Errorf("failed to update recipient broadcast status to failed: %v", errUpdate)
					return
				}
				continue
			}

			// Update recipient broadcast status
			_, err = h.service.UpdateRecipientBroadcastStatus(ctx, recipientID, broadcastID, "Success")
			if err != nil {
				ch <- err
				return
			}

			// Delay between sending messages
			time.Sleep(5 * time.Second)
		}

		// Close the channel once all processing is done
		close(ch)
	}()

	// Return after starting the goroutine
	return nil
}

func (h *BroadcastHandler) HandlePecatuBroadcast(c *fiber.Ctx) error {
	// Check if client is logged in
	if h.client.Store.ID == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Please login first",
		})
	}

	// Parse the payload
	payload := new(SendMessagePayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	// Get the broadcast ID from the payload
	broadcastID := payload.BroadcastID

	// Check if there are any recipients in the broadcast
	ctx := context.Background()
	hasRecipients, err := h.service.IsAnyRecipientInBroadcast(ctx, broadcastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to check recipients",
		})
	}
	if !hasRecipients {
		return c.Status(200).JSON(fiber.Map{
			"message": "No recipients to send broadcast",
		})
	}

	// Update broadcast status to "Starting"
	_, err = h.service.UpdateBroadcastStatus(ctx, broadcastID, "Starting")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update broadcast status",
		})
	}

	// Get the broadcast message
	broadcastMessage, err := h.service.GetBroadcastMessage(ctx, broadcastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get broadcast message",
		})
	}

	// Get recipients by broadcast ID
	recipients, err := h.service.GetAllRecipientByBroadcastID(ctx, broadcastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get recipients",
		})
	}

	// Send an immediate response
	response := fiber.Map{
		"message": "Message sending process has started",
	}
	if err := c.JSON(response); err != nil {
		return err // Handle any error in sending the response
	}

	// Start the background process
	go func() {
		ch := make(chan error, len(recipients["recipients"].([]map[string]interface{})))

		for _, recipient := range recipients["recipients"].([]map[string]interface{}) {
			recipientID, ok := recipient["id"].(int)
			if !ok {
				ch <- fmt.Errorf("recipient ID is not an integer")
				return
			}
			whatsappNumber, ok := recipient["whatsapp_number"].(string)
			if !ok {
				ch <- fmt.Errorf("WhatsApp number is not a string")
				return
			}

			recipientName, ok := recipient["recipient_name"].(string)
			if !ok {
				ch <- fmt.Errorf("recipient name is not a string")
				return
			}

			recipientIdentifier, ok := recipient["recipient_unique_identifier"].(string)
			if !ok {
				ch <- fmt.Errorf("identifier is not a string")
				return
			}

			// Ensure client is connected
			if !h.client.IsConnected() {
				err := h.client.Connect()
				if err != nil {
					ch <- fmt.Errorf("failed to connect to WhatsApp: %v", err)
					return
				}
			}

			// Prepare recipient JID and message
			recipientJID := types.JID{
				User:   whatsappNumber,
				Server: "s.whatsapp.net",
			}
			var msg *waProto.Message
			// Menggantikan teks {name} menggunakan recipientName dan {identifier} menggunakan recipientIdentifier untuk setiap pesan
			message := strings.ReplaceAll(broadcastMessage["broadcast_message"].(string), "{name}", recipientName)
			message = strings.ReplaceAll(message, "{identifier}", recipientIdentifier)

			if containsURL(message) {
				preview, err := fetchLinkPreview(message)
				if err == nil && preview != nil {
					msg = &waProto.Message{
						ExtendedTextMessage: &waProto.ExtendedTextMessage{
							Text:          proto.String(message),
							MatchedText:   proto.String(preview.MatchedText),
							CanonicalURL:  proto.String(preview.CanonicalURL),
							Title:         proto.String(preview.Title),
							Description:   proto.String(preview.Description),
							JPEGThumbnail: preview.JPEGThumbnail,
						},
					}
				}
			}
			if msg == nil {
				msg = &waProto.Message{
					Conversation: proto.String(message),
				}
			}

			// Mengirim pesan
			_, err = h.client.SendMessage(ctx, recipientJID, msg)
			if err != nil {
				// Jika pengiriman pesan gagal, lanjutkan ke penerima berikutnya dan update status menjadi gagal
				_, errUpdate := h.service.UpdateRecipientBroadcastStatus(ctx, recipientID, broadcastID, "Failed")
				if errUpdate != nil {
					ch <- fmt.Errorf("failed to update recipient broadcast status to failed: %v", errUpdate)
					return
				}
				continue
			}

			// Update recipient broadcast status
			_, err = h.service.UpdateRecipientBroadcastStatus(ctx, recipientID, broadcastID, "Success")
			if err != nil {
				ch <- err
				return
			}

			// Delay between sending messages
			time.Sleep(5 * time.Second)
		}

		// Close the channel once all processing is done
		close(ch)
	}()

	// Return after starting the goroutine
	return nil
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
