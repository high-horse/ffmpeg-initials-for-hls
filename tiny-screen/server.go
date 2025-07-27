package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tf "tiny-screen/ffmpeg-core"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func serve() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HLS server is running.\nUse /hls/sample-hls-op-m3u8 to access the stream.")
	})

	e.POST("/upload", submitMedia)
	e.GET("/hls/all", hlsMetadata)
	e.GET("/hls/*", hlsHandler)

	port := ":9090"
	fmt.Printf("🔊 Serving HLS on http://localhost%s\n", port)
	e.Logger.Fatal(e.Start(port))
}

func hlsHandler(c echo.Context) error {
	if c.Request().Method == http.MethodOptions {
		return c.NoContent(http.StatusOK)
	}

	filePath := c.Param("*")
	// filePath := filepath.Join("outputs", relativePath)
	// filePath = filepath.Clean(filePath)

	log.Printf("Serving file: %s", filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "File not found")
	}

	// Set appropriate Content-Type
	switch ext := filepath.Ext(filePath); ext {
	case ".m3u8":
		c.Response().Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	case ".ts":
		c.Response().Header().Set("Content-Type", "video/MP2T")
	default:
		return c.String(http.StatusUnsupportedMediaType, "Unsupported file type")
	}

	return c.File(filePath)
}

func hlsMetadata(c echo.Context) error {
	// Define path to central meta.json file
	metaFilePath := "meta.json"

	// Read the JSON file
	data, err := os.ReadFile(metaFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, return empty list
			return c.JSON(http.StatusOK, []map[string]string{})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read meta.json: " + err.Error()})
	}

	// Parse JSON content
	var registry []map[string]string
	if err := json.Unmarshal(data, &registry); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parse meta.json: " + err.Error()})
	}

	// Return metadata as JSON
	return c.JSON(http.StatusOK, registry)
}

func submitMedia(c echo.Context) error {
	file, err := c.FormFile("media_file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to get file: " + err.Error()})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to open file: " + err.Error()})
	}
	defer src.Close()

	// 2. Save to input-samples/
	inputPath := filepath.Join("input-samples", file.Filename)
	out, err := os.Create(inputPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to save file: " + err.Error()})
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to write file: " + err.Error()})
	}
	// 3. Encode HLS
	x11 := tf.NewFfmpeg()
	hlsPath, err := x11.HlsEncodeLocal(inputPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to encode HLS: " + err.Error()})
	}

	// 4. Prepare metadata
	outputDir := filepath.Dir(hlsPath)
	id := filepath.Base(outputDir) // e.g. test_1721998123

	newEntry := map[string]string{
		"id":   id,
		"name": file.Filename,
		"path": hlsPath,
	}

	// 5. Write metadata to JSON
	metaFilePath := "meta.json"
	var registry []map[string]string

	if data, err := os.ReadFile(metaFilePath); err == nil && len(data) > 0 {
		json.Unmarshal(data, &registry) // ignore unmarshal errors silently
	}

	// 6. Append new entry and save
	registry = append(registry, newEntry)

	metaFile, err := os.Create(metaFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update meta.json: " + err.Error()})
	}
	defer metaFile.Close()

	encoder := json.NewEncoder(metaFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(registry); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to write meta.json: " + err.Error()})
	}

	// 7. Done
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Media uploaded and processed",
		"meta":    newEntry,
	})
}

func serve_() {
	http.HandleFunc("/hls/", hlsHandler_)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HLS server is running.\nUse /hls/sample-hls-op-m3u8 to access the stream."))
	})

	port := ":9090"
	fmt.Printf("🔊 Serving HLS on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func hlsHandler_(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Build and clean file path
	relativePath := strings.TrimPrefix(r.URL.Path, "/hls")
	filePath := filepath.Join("outputs", "hls", relativePath)
	filePath = filepath.Clean(filePath)
	log.Printf("Requested filePath: %s", filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Log the extension explicitly
	ext := filepath.Ext(filePath)
	log.Printf("File extension: %q", ext) // Use %q to show quoted string, helps detect hidden characters

	// Set appropriate Content-Type headers
	switch ext {
	case ".m3u8":
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	case ".ts":
		w.Header().Set("Content-Type", "video/MP2T")
	default:
		log.Printf("Unsupported file extension %q for: %s", ext, filePath)
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}
