package hassio

import "fmt"

// https://www.home-assistant.io/docs/mqtt/discovery/

type sensorModel struct {
	Name              string        `json:"name"` 				 // Sensor Name
	AvailabilityTopic string        `json:"availability_topic"`  // The MQTT topic subscribed to receive availability (online/offline) updates.
	PayloadAvailable  string        `json:"payload_available"`   // The payload that represents the available state.
	PayloadNotAvailable string      `json:"payload_not_available"` // The payload that represents the unavailable state.
	StateTopic        string        `json:"state_topic"`  		 // State Topic
	Device            deviceModel   `json:"device"`              // Information about the device this sensor is a part of to tie it into the device
	DeviceClass       interface{}   `json:"device_class,omitempty"`        // The type/class of the sensor to set the icon in the frontend.
	UnitOfMeasurement interface{}   `json:"unit_of_measurement,omitempty"` // °C
	ValueTemplate     string        `json:"value_template"`      // "{{ value_json.temperature}}"
	UniqueId          string        `json:"unique_id"`           // devicename???
	Icon              string        `json:"icon"`                // "mdi:thermometer"
	Qos               int           `json:"qos"`                 // 1
}

type Sensor struct {
	device 				*Device
	id 					string
	model      			sensorModel
	sensorRegistered 	bool
	sensorValueFunction SensorValueFunction
}

type SensorValueFunction func() interface{}

// Add Sensor To Device
func (device *Device) AddSensor(
	name string,
	id string,
	deviceClass interface{},
	unitOfMeasurement interface{},
	valueTemplateFilter string,
	icon string,
	sensorValueFunction SensorValueFunction) Sensor{

	sensor := Sensor{
		device: device,
		id: id,
		model: sensorModel{
			Name: name,
			StateTopic:        device.sensorStateTopic,
			AvailabilityTopic: device.client.lwtTopic,
			PayloadAvailable:  "online",
			PayloadNotAvailable: "offline",
			Device:            device.model,
			DeviceClass:       deviceClass,
			UnitOfMeasurement: unitOfMeasurement,
			ValueTemplate:     fmt.Sprintf("{{ value_json.%s %s}}", id, valueTemplateFilter),
			UniqueId:          fmt.Sprintf("sensor_%s_%s", id, device.client.id),
			Icon:              icon,
			Qos:               1,
		},
		sensorRegistered: false,
		sensorValueFunction: sensorValueFunction,
	}


	device.sensorValues[id] = sensor.sensorValueFunction()
	device.sensors = append(device.sensors, sensor)
	device.sensorsChanged = true

	return sensor
}

// Update Sensor Value
func (sensor *Sensor) UpdateSensorValue(){
	value := sensor.sensorValueFunction()
	sensor.SetValue(value)
}

// Set Sensor Value
func (sensor *Sensor) SetValue(value interface{}){
	sensor.device.sensorValues[sensor.id] = value
	sensor.device.sensorValuesChanged = true
}