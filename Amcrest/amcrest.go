package amcrest

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"

	//"io"
	"log"
	"net/http"

	//"os"
	"strings"
)

type Camera struct {
	URI      string
	Client   *http.Client
	Username string
	Password string
}

func Init(uri string, username string, password string) *Camera {

	// proxyURL, err := url.Parse("http://127.0.0.1:8080")

	// if err != nil {
	// 	log.Fatalf("Invalid proxy URL: %v", err)
	// }

	transport := &http.Transport{
	//	Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	// ********** CREATE HTTP CLIENT ************

	client := &http.Client{
		Transport: transport,
	}

	return &Camera{URI: uri, Username: username, Password: password, Client: client}
}

func (cam *Camera) SetNounceCount(nc string) {
	// cam.NounceCount = nc
}

func (cam *Camera) LoadAuth() (string, error) {

	// Initial request to get WWW-Authenticate header

	req, err := http.NewRequest("GET", cam.URI, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Perform the initial request
	res, err := cam.Client.Do(req)
	if err != nil {
		log.Fatalf("Error making initial request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	// Check for 401 Unauthorized
	if res.StatusCode != http.StatusUnauthorized {
		return "", fmt.Errorf("expected 401 Unauthorized, got %v", res.Status)

	}

	// Parse the WWW-Authenticate header
	authHeader := res.Header.Get("WWW-Authenticate")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Digest") {
		return "", fmt.Errorf("digest authentication not supported or invalid WWW-Authenticate header: %s", authHeader)
	}
	params := parseAuthHeader(authHeader)

	// Extract parameters
	realm := params["realm"]
	nonce := params["nonce"]
	qop := params["qop"]
	opaque := params["opaque"] // Ensure opaque is extracted and non-empty
	nc := "00000001"
	cnonce := generateCnonce() // Generate a unique cnonce for every request

	method := "GET"
	uri := "/cgi-bin/snapshot.cgi"

	// Compute HA1, HA2, and response hashes
	ha1 := md5Hash(fmt.Sprintf("%s:%s:%s", cam.Username, realm, cam.Password))
	ha2 := md5Hash(fmt.Sprintf("%s:%s", method, uri))
	response := md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, nonce, nc, cnonce, qop, ha2))

	// Build the Authorization header
	return fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s", qop=%s, nc=%s, cnonce="%s", opaque="%s"`,
		cam.Username, realm, nonce, uri, response, qop, nc, cnonce, opaque,
	), nil

}

func (cam *Camera) GetSnapshot() ([]byte, error) {

	auth, err := cam.LoadAuth()

	if err != nil {
		log.Fatal(err)
	}

	// Second request with Digest Authorization header
	req, err := http.NewRequest("GET", cam.URI, nil)
	if err != nil {
		//return nil, fmt.Errorf("error creating second request: %v", err)
	}
	req.Header.Set("Authorization", auth)

	// Perform the authenticated request
	res, err := cam.Client.Do(req)
	if err != nil {
		//return nil, fmt.Errorf("error making authenticated request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil

}

// parseAuthHeader parses the WWW-Authenticate header into a map of key-value pairs
func parseAuthHeader(header string) map[string]string {
	params := make(map[string]string)

	// Remove the "Digest " prefix if present
	header = strings.TrimPrefix(header, "Digest ")

	// Split the header by commas
	fields := strings.Split(header, ",")
	for _, field := range fields {
		// Split each field into key=value
		parts := strings.SplitN(strings.TrimSpace(field), "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), `"`) // Remove quotes around values
			params[key] = value
		}
	}
	return params
}

// md5Hash computes the MD5 hash of a string
func md5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// generateCnonce generates a random client nonce
func generateCnonce() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
