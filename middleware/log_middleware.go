package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	logDir := "/var/log/app"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Set proper ownership if running as root
	if os.Geteuid() == 0 {
		// You might need to adjust these IDs based on your container user
		if err := os.Chown(logDir, 1000, 1000); err != nil {
			log.Warnf("Failed to change log directory ownership: %v", err)
		}
	}

	// Create log file with timestamp and proper permissions
	logFile := filepath.Join(logDir, "application.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Ensure proper file ownership
	if os.Geteuid() == 0 {
		if err := file.Chown(1000, 1000); err != nil {
			log.Warnf("Failed to change log file ownership: %v", err)
		}
	}

	// Configure logrus
	log.SetOutput(file)
	log.SetLevel(logrus.InfoLevel)

}

func LoggingMiddleware(serviceName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		var requestBody []byte
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Error("Error reading request body: ", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			requestBody = body

			// Restore the request body so it can be read again
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		// Create a response writer to capture the response
		responseWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		responseBuffer := new(bytes.Buffer)

		// Call the next handler
		next.ServeHTTP(responseWriter, r)

		responseBuffer.Write(responseWriter.body)

		// Log the response details
		log.WithFields(logrus.Fields{
			"timestamp":   time.Now().Format(time.RFC3339),
			"level":       "INFO",
			"service":     serviceName,
			"environment": "production",
			"message":     "Handled request",
			"request": map[string]interface{}{
				"method":  r.Method,
				"url":     r.URL.String(),
				"headers": r.Header,
				"body":    string(requestBody), // Log the request body
			},
			"response": map[string]interface{}{
				"statusCode": responseWriter.statusCode,
				"body":       responseBuffer.String(),
			},
		}).Info("Response logged")
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// WriteHeader captures the status code
func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	w.body = b // Capture the response body
	return w.ResponseWriter.Write(b)
}
