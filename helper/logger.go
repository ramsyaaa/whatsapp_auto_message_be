package helper

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type LoggerConfig struct {
	Format       string
	TimeFormat   string
	TimeZone     string
	Output       *os.File
	CustomTags   map[string]logger.LogFunc
	Done         func(*fiber.Ctx, []byte)
	DisableColor bool
}

func NewLogger(config LoggerConfig) fiber.Handler {
	return logger.New(logger.Config{
		Format:     config.Format,
		TimeFormat: config.TimeFormat,
		TimeZone:   config.TimeZone,
		Output:     config.Output,
		CustomTags: config.CustomTags,
		Done:       config.Done,
	})
}

func DefaultLogger() fiber.Handler {
	return NewLogger(LoggerConfig{
		Format: "[${time}] ${ip} ${status} - ${method} ${path}\n",
	})
}

func RequestIDLogger() fiber.Handler {
	return requestid.New()
}

func CustomLogger(config LoggerConfig) fiber.Handler {
	return NewLogger(config)
}

func CustomFileLogger(filePath string) (*os.File, fiber.Handler, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return file, NewLogger(LoggerConfig{
		Output: file,
	}), err
}

func CustomTagLogger(customTags map[string]logger.LogFunc) fiber.Handler {
	return NewLogger(LoggerConfig{
		CustomTags: customTags,
	})
}

func CallbackLogger(done func(*fiber.Ctx, []byte)) fiber.Handler {
	return NewLogger(LoggerConfig{
		Done: done,
	})
}

func DisableColorLogger() fiber.Handler {
	return NewLogger(LoggerConfig{
		DisableColor: true,
	})
}

func LogToFile() fiber.Handler {
	// Tentukan jalur file log berdasarkan tanggal saat ini
	filePath := "logs/" + time.Now().Format("20060102") + "-log.log"

	// Cek apakah file log untuk tanggal saat ini sudah ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Jika file tidak ada, buat file baru
		if err := os.MkdirAll(filepath.Dir(filePath), 0666); err != nil {
			log.Fatalf("error creating directory: %v", err)
		}
		if _, err := os.Create(filePath); err != nil {
			log.Fatalf("error creating file: %v", err)
		}
	}

	// Buka file dengan mode baca/tulis, buat atau tambah
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// Kembalikan middleware logger fiber, mengarahkannya ke penulis file
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Local",
		Output:     file, // Tidak ada lagi stdout karena perubahan untuk tidak mencetak log di command line
	})
}
