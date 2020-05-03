package hassio

import (
	"encoding/json"
	"fmt"
	"log"
)

// https://www.home-assistant.io/docs/mqtt/discovery/

type deviceModel struct {
	Name         string `json:"name"` 			// Human Readable Device Name
	Model        string `json:"model"` 			// Board Model
	Manufacturer string `json:"manufacturer"` 	// Board Manufacturer
	SwVersion    string `json:"sw_version"` 	// Software Version
	Identifiers  []string `json:"identifiers"`	// Unique Device Identifiers (Serial, MAC, etc)
}

type Device struct {
	client Client
	changed    			bool
	//id 		   			string
	model      			deviceModel
	sensors	   			[]Sensor
	sensorStateTopic 	string
	sensorsChanged 		bool
	sensorValues 		map[string]interface{}
	sensorValuesChanged bool
}

// Add Device to Client
func (client Client) AddDevice(name string, model string, vendor string, swVersion string, identifiers []string) Device{
	device := Device{
		client: client,
		changed:    true,
		model: deviceModel{
			Name: name,
			Model: model,
			Manufacturer: vendor,
			SwVersion: swVersion,
			Identifiers: identifiers,
		},
		sensorStateTopic: fmt.Sprintf("%s/state", client.baseTopic),
		sensorsChanged: false,
		sensorValues: map[string]interface{}{},
		sensorValuesChanged: false,
	}
	return device
}

// Update Sensor Values
func (device *Device) UpdateSensorValues() {
	for index, _ := range device.sensors {
		device.sensors[index].UpdateSensorValue()
	}
}

// Submit Changes to MQTT Broker
func (device *Device) SubmitChanges() {
	// Submit Sensor Values
	if device.sensorValuesChanged == true {
		payload, err := json.Marshal(device.sensorValues)
		if err != nil{
			panic(err)
		}

		err = device.client.Publish(device.sensorStateTopic, 1, false, string(payload))
		if err != nil {
			panic(err)
		}

		device.sensorValuesChanged = false
	}

	// Register Sensors in Home Assistant
	if device.sensorsChanged {
		for index, _ := range device.sensors {
			if !device.sensors[index].sensorRegistered {
				log.Println("Registering sensor:", device.sensors[index].model.Name)

				payload, err := json.Marshal(device.sensors[index].model)
				if err != nil {
					panic(err)
				}

				topic := fmt.Sprintf("homeassistant/sensor/%s/%s_%s/config", device.client.id, device.client.id, device.sensors[index].id)
				err = device.client.Publish(topic, 1, true, string(payload))
				if err != nil {
					panic(err)
				}

				device.sensors[index].sensorRegistered = true
			}
		}
	}
}