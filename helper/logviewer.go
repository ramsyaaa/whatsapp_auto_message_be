package helper

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Route to get available log files
func GetLogFiles(c *fiber.Ctx) error {
	// Get all log files in the logs directory
	files, err := ioutil.ReadDir("./logs")
	if err != nil {
		return c.Status(500).SendString("Unable to read log directory")
	}

	var logFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") { // Ensure we're only listing `.log` files
			logFiles = append(logFiles, file.Name())
		}
	}

	// Return the list of log files as JSON
	return c.JSON(fiber.Map{
		"files": logFiles,
	})
}

// Route to get the logs from a selected log file
func GetLogFileContent(c *fiber.Ctx) error {
	filename := c.Params("filename")
	logFilePath := filepath.Join("./logs", filename)

	// Read the log file content
	content, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		return c.Status(500).SendString("Unable to read log file")
	}

	// Parse the log content
	logEntries := parseLogContent(string(content))

	// Return parsed log entries as JSON
	return c.JSON(fiber.Map{
		"logs": logEntries,
	})
}

// A simple log parser that splits on the `|` delimiter
func parseLogContent(content string) []map[string]string {
	var logs []map[string]string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Split the log line based on the `|` separator
		parts := strings.Split(line, "|")
		if len(parts) < 7 {
			continue
		}

		// Trim any leading or trailing spaces from each part
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		// Ensure we have the correct log structure
		logEntry := map[string]string{
			"timestamp":       parts[0], // timestamp (e.g., 2024-11-11T21:05:21+07:00)
			"status_code":     parts[1], // status code (e.g., 200)
			"latency":         parts[2], // latency (e.g., 3.693ms)
			"ip":              parts[3], // IP (e.g., 127.0.0.1)
			"method":          parts[4], // HTTP method (e.g., POST)
			"path":            parts[5], // Path (e.g., /api/v1/geomapping/sensor-list)
			"additional_info": parts[6], // Additional info (e.g., "-")
		}

		logs = append(logs, logEntry)
	}

	return logs
}
