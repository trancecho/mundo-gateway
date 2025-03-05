package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/po"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	// Initialize the gateway
	domain.GatewayGlobal = domain.NewGateway()

	// Insert mock data for the test
	mockPrefix := po.Prefix{Name: "/api/v1", ServiceId: 1}
	log.Println("Inserting mockPrefix:", mockPrefix)
	if err := domain.GatewayGlobal.DB.Create(&mockPrefix).Error; err != nil {
		log.Printf("Error inserting mockPrefix: %v\n", err)
	} else {
		log.Println("Successfully inserted mockPrefix")
	}
	domain.GatewayGlobal.Prefixes = append(domain.GatewayGlobal.Prefixes, mockPrefix)

	mockService := po.Service{
		Name:   "MockService",
		Prefix: "/api/v1",
		APIs: []po.API{
			{
				HttpPath:   "/ping",
				HttpMethod: http.MethodGet,
			},
		},
	}
	log.Println("Inserting mockService:", mockService)
	if err := domain.GatewayGlobal.DB.Create(&mockService).Error; err != nil {
		log.Printf("Error inserting mockService: %v\n", err)
	} else {
		log.Println("Successfully inserted mockService")
	}
	domain.GatewayGlobal.Services = append(domain.GatewayGlobal.Services, mockService)
	// Start the mock service for testing
	startMockService()
}

// mock target service running on port 8081
func startMockService() {
	// Initialize a Gin engine for the mock service
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Define the /ping route for the mock service
	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Start the mock service on port 8081
	go func() {
		if err := r.Run(":8081"); err != nil {
			log.Fatalf("Could not start mock service on port 8081: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)
}

// TestReverseProxy 测试反向代理是否能正确工作
func TestReverseProxy(t *testing.T) {
	// Initialize a Gin engine for testing
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Any("/*path", controller.HandleRequestController)

	// Log the initialization of the server
	log.Println("Initializing the server for testing")

	// Mock request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w := httptest.NewRecorder()

	// Log the request details
	log.Printf("Sending request: %s %s\n", req.Method, req.URL.Path)

	// Execute the request
	r.ServeHTTP(w, req)

	// Log the response status code and body
	log.Printf("Response Status Code: %d\n", w.Code)
	log.Printf("Response Body: %s\n", w.Body.String())

	// Check the status code
	assert.Equal(t, 200, w.Code)

	// Check if the response is what we expect (assuming your service is correctly set up)
	expectedResponse := `{"message":"pong"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
