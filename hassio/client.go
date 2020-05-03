package hassio

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

type Client struct {
	client mqtt.Client
	lwtTopic string
	baseTopic string
	id string
}

// Create New Client To Broker
func NewClient(broker string, id string, user string, password string) Client {
	baseTopic := fmt.Sprintf("system-telemetry/%s", id)
	lwtTopic := fmt.Sprintf("%s/lwt", baseTopic)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(true)
	opts.SetWill(lwtTopic, "offline", 1, true)

	client := Client{
		client: mqtt.NewClient(opts),
		lwtTopic: lwtTopic,
		baseTopic: baseTopic,
		id: id,
	}

	return client
}

// Connect To Broker
func (client Client) Connect() {
	Connect:
	log.Print("Attempting to connect to broker.. ")

	timeout := 0 * time.Second
	for {
		token := client.client.Connect()
		token.Wait()
		err := token.Error()
		if err == nil {
			break
		}
		log.Println("ERROR: Failed with exception:", err.Error())

		// Retry Timeout
		if timeout < (30 * time.Second) {
			timeout = timeout + (5 * time.Second)
		}
		log.Println("Waiting for", timeout, " seconds before retying")
		time.Sleep(timeout)
		log.Print("Attempting to connect to broker.. ")
	}

	log.Println("Connected")

	// Send Last Will Testament message
	err := client.Publish(client.lwtTopic, 1, true, "online")
	if err != nil && !client.client.IsConnected() {
		log.Println("ERROR: Failed to register Last Will Testament message due to:", err)

		if !client.client.IsConnected() {
			goto Connect
		}
	}
}

// Disconnect From Broker
func (client Client) Disconnect(){
	client.Publish(client.lwtTopic, 1, true, "offline")
	client.client.Disconnect(250)
	// TODO: Add semaphore to avoid connecting while disconnecting?
}

// Submit Message
func (client Client) Publish(topic string, qos byte, retained bool, payload string) error{
	Publish:

	token := client.client.Publish(topic, qos, retained, payload)
	token.Wait()
	err := token.Error()

	if err != nil {
		if !client.client.IsConnected() {
			client.Connect()
			goto Publish
		} else {
			return err
		}
	}

	return nil
}

