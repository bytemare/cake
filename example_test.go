package cake

import (
	"bytes"
	"fmt"
)

func ExampleCake() {
	serverID := []byte("server")
	username := []byte("client")
	secret := []byte("password")

	// Set up client
	client, err := Client(serverID, username, secret)
	if err != nil {
		panic(err)
	}

	// Set up server
	server, err := Server(serverID, username, secret)
	if err != nil {
		panic(err)
	}

	// Client starts protocol and sends message to server
	message1, err := client.Start()
	if err != nil {
		panic(err)
	}

	// The server awaits and interprets the client's message,
	// and returns the server response to the client
	message2, err := server.AuthenticateClient(message1)
	if err != nil {
		panic(err)
	}

	// The client awaits an interprets the server's response
	err = client.Finish(message2)
	if err != nil {
		panic(err)
	}

	// The protocol is finished, and both parties now share the same secret session key
	if bytes.Equal(server.SessionKey(), client.SessionKey()) {
		fmt.Println("Success ! Both parties share the same secret session key !")
	} else {
		fmt.Printf("Failed. Client and server keys are different.")
	}
	// Output: Success ! Both parties share the same secret session key !
}
