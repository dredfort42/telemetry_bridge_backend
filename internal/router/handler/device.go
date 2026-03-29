package handler

import (
	"net/http"
	"telemetry_bridge/internal/router/model"

	log "github.com/dredfort42/go_logger"
	"github.com/gin-gonic/gin"
)

// RegisterDevice handles the registration of a new device. It expects a JSON payload with device details and responds with the registration status.
func RegisterDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var device model.DeviceRegisterRequest

		if err := c.ShouldBindJSON(&device); err != nil {
			status := http.StatusBadRequest
			response := model.GetErrorResponse(status, "The request body is invalid or missing required fields: "+err.Error(), c.Request.RequestURI, nil)

			c.AbortWithStatusJSON(status, response)

			return
		}

		// device, err := db.RegisterDevice(req)
		// if err != nil {
		// 	if db.IsDeviceNameAlreadyExistsErr(err) {
		// 		status := http.StatusConflict
		// 		response := model.GetErrorResponse(status, "The device name is already registered", c.Request.RequestURI, nil)

		// 		c.AbortWithStatusJSON(status, response)

		// 		return
		// 	}

		// 	status := http.StatusInternalServerError
		// 	response := model.GetErrorResponse(status, "Failed to register device: "+err.Error(), c.Request.RequestURI, nil)

		// 	c.AbortWithStatusJSON(status, response)

		// 	return
		// }

		// response := model.GetSuccessResponse(http.StatusOK, "Device registered successfully", c.Request.RequestURI, device)

		log.Info.Printf("Registered device: %v", device)

		c.JSON(http.StatusOK, gin.H{"status": "registered", "data": device})

		// body := make(map[string]any)
		// if err := c.BindJSON(&body); err != nil {
		// 	c.JSON(400, gin.H{"error": "invalid JSON"})
		// 	return
		// }

		// log.Debug.Printf("Received registration: %v", body)

		// c.JSON(200, gin.H{"status": "registered", "data": body})
	}
}
