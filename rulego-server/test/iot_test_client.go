package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// IoTSensorData represents the IoT sensor data structure
type IoTSensorData struct {
	DeviceID  string    `json:"deviceId"`
	Timestamp string    `json:"timestamp"`
	Data      SensorData `json:"data"`
}

// SensorData represents the sensor measurements
type SensorData struct {
	Temperature *float64 `json:"temperature,omitempty"`
	Humidity    *float64 `json:"humidity,omitempty"`
	Pressure    *float64 `json:"pressure,omitempty"`
}

// TestClient represents the MQTT test client
type TestClient struct {
	client mqtt.Client
	topic  string
}

// NewTestClient creates a new MQTT test client
func NewTestClient(brokerURL, clientID, topic string) (*TestClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)

	// Connection callback
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("Connected to MQTT broker: %s", brokerURL)
	})

	// Connection lost callback
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	return &TestClient{
		client: client,
		topic:  topic,
	}, nil
}

// Close closes the MQTT client connection
func (tc *TestClient) Close() {
	tc.client.Disconnect(250)
}

// PublishSensorData publishes sensor data to MQTT topic
func (tc *TestClient) PublishSensorData(data IoTSensorData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal sensor data: %v", err)
	}

	token := tc.client.Publish(tc.topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %v", token.Error())
	}

	log.Printf("Published sensor data: %s", string(payload))
	return nil
}

// GenerateRandomSensorData generates random sensor data for testing
func GenerateRandomSensorData(deviceID string) IoTSensorData {
	data := IoTSensorData{
		DeviceID:  deviceID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      SensorData{},
	}

	// Randomly generate different types of sensor data
	switch rand.Intn(4) {
	case 0:
		// Temperature only
		temp := 15.0 + rand.Float64()*50.0 // 15-65°C
		data.Data.Temperature = &temp
	case 1:
		// Humidity only
		humidity := 20.0 + rand.Float64()*60.0 // 20-80%
		data.Data.Humidity = &humidity
	case 2:
		// Pressure only
		pressure := 950.0 + rand.Float64()*100.0 // 950-1050 hPa
		data.Data.Pressure = &pressure
	case 3:
		// All sensors
		temp := 15.0 + rand.Float64()*50.0
		humidity := 20.0 + rand.Float64()*60.0
		pressure := 950.0 + rand.Float64()*100.0
		data.Data.Temperature = &temp
		data.Data.Humidity = &humidity
		data.Data.Pressure = &pressure
	}

	return data
}

// GenerateAbnormalSensorData generates abnormal sensor data for testing alerts
func GenerateAbnormalSensorData(deviceID string, alertLevel string) IoTSensorData {
	data := IoTSensorData{
		DeviceID:  deviceID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      SensorData{},
	}

	switch alertLevel {
	case "warning":
		// Generate warning level data
		temp := 36.0 + rand.Float64()*8.0 // 36-44°C (warning range)
		humidity := 81.0 + rand.Float64()*8.0 // 81-89% (warning range)
		data.Data.Temperature = &temp
		data.Data.Humidity = &humidity
	case "critical":
		// Generate critical level data
		temp := 46.0 + rand.Float64()*14.0 // 46-60°C (critical range)
		humidity := 91.0 + rand.Float64()*8.0 // 91-99% (critical range)
		data.Data.Temperature = &temp
		data.Data.Humidity = &humidity
	}

	pressure := 950.0 + rand.Float64()*100.0 // Normal pressure
	data.Data.Pressure = &pressure

	return data
}

// RunTestScenario runs different test scenarios
func RunTestScenario(client *TestClient, scenario string, duration time.Duration) {
	log.Printf("Starting test scenario: %s for %v", scenario, duration)
	
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	timeout := time.After(duration)
	devices := []string{"sensor_001", "sensor_002", "sensor_003"}
	
	for {
		select {
		case <-timeout:
			log.Printf("Test scenario %s completed", scenario)
			return
		case <-ticker.C:
			deviceID := devices[rand.Intn(len(devices))]
			
			var data IoTSensorData
			switch scenario {
			case "normal":
				data = GenerateRandomSensorData(deviceID)
			case "warning":
				data = GenerateAbnormalSensorData(deviceID, "warning")
			case "critical":
				data = GenerateAbnormalSensorData(deviceID, "critical")
			case "mixed":
				// Mix of normal and abnormal data
				switch rand.Intn(4) {
				case 0:
					data = GenerateAbnormalSensorData(deviceID, "warning")
				case 1:
					data = GenerateAbnormalSensorData(deviceID, "critical")
				default:
					data = GenerateRandomSensorData(deviceID)
				}
			}
			
			if err := client.PublishSensorData(data); err != nil {
				log.Printf("Error publishing data: %v", err)
			}
		}
	}
}

// PublishInvalidData publishes invalid data to test filtering
func PublishInvalidData(client *TestClient) {
	log.Println("Publishing invalid data to test filtering...")
	
	// Invalid data scenarios
	invalidDataList := []map[string]interface{}{
		// Missing deviceId
		{
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"temperature": 25.0,
			},
		},
		// Missing timestamp
		{
			"deviceId": "sensor_001",
			"data": map[string]interface{}{
				"temperature": 25.0,
			},
		},
		// Missing data
		{
			"deviceId": "sensor_001",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		// Invalid temperature range
		{
			"deviceId": "sensor_001",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"temperature": -100.0, // Out of range
			},
		},
		// Invalid humidity range
		{
			"deviceId": "sensor_001",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"humidity": 150.0, // Out of range
			},
		},
		// Invalid pressure range
		{
			"deviceId": "sensor_001",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"pressure": 500.0, // Out of range
			},
		},
	}
	
	for i, invalidData := range invalidDataList {
		payload, _ := json.Marshal(invalidData)
		token := client.client.Publish(client.topic, 1, false, payload)
		if token.Wait() && token.Error() != nil {
			log.Printf("Error publishing invalid data %d: %v", i+1, token.Error())
		} else {
			log.Printf("Published invalid data %d: %s", i+1, string(payload))
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// MQTT configuration
	brokerURL := "tcp://localhost:1883"
	clientID := "iot_test_client"
	topic := "iot/sensors/data"
	
	// Create MQTT client
	client, err := NewTestClient(brokerURL, clientID, topic)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	defer client.Close()
	
	log.Println("IoT Test Client started")
	log.Println("Available commands:")
	log.Println("  1. normal - Send normal sensor data")
	log.Println("  2. warning - Send warning level data")
	log.Println("  3. critical - Send critical level data")
	log.Println("  4. mixed - Send mixed normal/abnormal data")
	log.Println("  5. invalid - Send invalid data to test filtering")
	log.Println("  6. continuous - Run continuous mixed data simulation")
	log.Println("  7. exit - Exit the program")
	
	// Interactive mode
	for {
		fmt.Print("\nEnter command (1-7): ")
		var choice string
		fmt.Scanln(&choice)
		
		switch choice {
		case "1", "normal":
			RunTestScenario(client, "normal", 30*time.Second)
		case "2", "warning":
			RunTestScenario(client, "warning", 30*time.Second)
		case "3", "critical":
			RunTestScenario(client, "critical", 30*time.Second)
		case "4", "mixed":
			RunTestScenario(client, "mixed", 60*time.Second)
		case "5", "invalid":
			PublishInvalidData(client)
		case "6", "continuous":
			log.Println("Starting continuous simulation (press Ctrl+C to stop)...")
			RunTestScenario(client, "mixed", 24*time.Hour) // Run for a very long time
		case "7", "exit":
			log.Println("Exiting...")
			return
		default:
			log.Println("Invalid choice. Please enter 1-7.")
		}
	}
} 