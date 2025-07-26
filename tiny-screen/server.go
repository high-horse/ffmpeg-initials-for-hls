package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func serve() {
	http.HandleFunc("/hls/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HLS server is running.\nUse /hls/sample-hls-op-m3u8 to access the stream."))
	})

	port := ":9090"
	fmt.Printf("🔊 Serving HLS on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func serve_0() {
	http.HandleFunc("/hls/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Build path to the requested file
		filePath := "outputs/hls" + strings.TrimPrefix(r.URL.Path, "/hls")

		// Set appropriate Content-Type headers
		switch filepath.Ext(filePath) {
		case ".m3u8":
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		case ".ts":
			w.Header().Set("Content-Type", "video/MP2T")
		default:
			http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
			return
		}

		// Serve the file
		http.ServeFile(w, r, filePath)
	})

	// Optionally serve a health check or root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HLS server is running.\nUse /hls/sample-hls-op-m3u8 to access the stream."))
	})

	// Start the server
	port := ":9090"
	fmt.Printf("🔊 Serving HLS on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
