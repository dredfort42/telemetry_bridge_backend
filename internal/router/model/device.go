package model

type DeviceRegisterRequest struct {
	DeviceInfo   DeviceInfo         `json:"device_info"`  // Basic information about the device
	Capabilities DeviceCapabilities `json:"capabilities"` // What the device can do (sensors, actuators, etc.)
}

type DeviceInfo struct {
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Firmware string `json:"firmware"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	MAC      string `json:"mac"`
}

type DeviceCapabilities struct {
	Sensors   *[]DeviceSensor   `json:"sensors,omitempty"`   // List of sensors the device has, if any
	Actuators *[]DeviceActuator `json:"actuators,omitempty"` // List of actuators the device has, if any
}

type DeviceSensor struct {
	ID            string     `json:"id"`             // Unique identifier for the sensor
	Type          string     `json:"type"`           // e.g., "temperature", "humidity", "pressure"
	Unit          string     `json:"unit"`           // e.g., "°C", "%", "hPa"
	Range         [2]float64 `json:"range"`          // [min, max] values the sensor can measure
	Resolution    float64    `json:"resolution"`     // Smallest change the sensor can detect
	MinDelay      int        `json:"min_delay"`      // Minimum delay between measurements in milliseconds
	ReadOnly      bool       `json:"read_only"`      // Indicates if the sensor is read-only or can be configured
	SamplingModes []string   `json:"sampling_modes"` // e.g., ["push", "pull"]
}

type DeviceActuator struct {
	ID       string                `json:"id"`               // Unique identifier for the actuator
	Type     string                `json:"type"`             // e.g., "switch", "motor", "valve"
	Commands []string              `json:"commands"`         // e.g., ["on", "off"] for a switch, ["start", "stop", "set_speed"] for a motor
	State    []string              `json:"state,omitempty"`  // e.g., ["on", "off"]
	Params   *DeviceActuatorParams `json:"params,omitempty"` // Additional parameters for actuators, e.g., speed range for motors
}

type DeviceActuatorParams struct {
	Speed *DeviceActuatorSpeedParam `json:"speed,omitempty"` // Parameters related to speed control, if applicable (e.g., for motors)
}

type DeviceActuatorSpeedParam struct {
	Min  int    `json:"min"`  // Minimum speed value
	Max  int    `json:"max"`  // Maximum speed value
	Unit string `json:"unit"` // Unit of speed, e.g., "RPM", "m/s"
}
