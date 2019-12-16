package main

func main() {
	// TODO: write a client with go-mqtt/client, it should do...
	// 1. connect to MQTT server
	// 2. subscribe "#" topic.
	// 3. publish a message "Hi MQTT Server" to "a/b/c/response" topic
	// 4. log messages on each phases

	done := make(chan struct{})
	<-done
}
