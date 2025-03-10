package isolatedfunctions

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func VaultStatusCheck() bool {
	// address := "127.0.0.1:8200"
	address := "host.docker.internal:8200" // Only should be uncommented when building for docker
	timeout := 2 * time.Second             // Timeout after 2 seconds

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}

	conn.Close() // Close connection if successful
	return true
}

func CreatePayload(payload map[string]string) string {
	var apiPayload []string

	for key, value := range payload {
		apiPayload = append(apiPayload, key+"="+value)
	}

	return strings.Join(apiPayload, "&")
}

func POSTforFormPayload(url string, payload url.Values) (response *http.Response, err error) {

	resp, err := http.PostForm(url, payload)
	if err != nil {
		log.Printf("Failed to fetch auth token from SNOW instance due to %q", err)
		return nil, err
	}

	return resp, err
}

func POSTjsonPayload(c *gin.Context, authToken, requestType, url string, payload []byte) (responseBody []byte, err error) {

	request, _ := http.NewRequestWithContext(c, requestType, url, bytes.NewBuffer(payload))

	request.Header.Add("Authorization", "Bearer "+authToken)
	request.Header.Set("Content-Type", "application/json") // Adjust content type as needed

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Unable to call update incident API", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseCode := resp.StatusCode

	if responseCode > 210 {
		log.Println("API call failed due to ", resp.Status)
		return nil, errors.New(resp.Status)
	}

	body, _ := io.ReadAll(resp.Body) // Reads the body from http.Response and converts it into []byte
	return body, err
}
