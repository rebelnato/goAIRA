package authorizations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/rebelnato/goAIRA/endpoints"
	"github.com/rebelnato/goAIRA/isolatedfunctions"
	"github.com/rebelnato/goAIRA/vault"
)

var snowAuthEndpoint string = endpoints.ConfigData.Endpoints["servicenow"]["base"].(string) + "oauth_token.do"

type SNOWResponse struct {
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken     string `json:"refresh_token"`
	RefreshTokenTime string `json:refresh_time`
}

func GetSNOWRefreshToken() (string, bool) {
	vaultData := vault.ReadSecrets("SNOW_refresh", "secret")
	currTime := time.Now().Unix()

	var refreshToken RefreshToken
	validRefreshToken := true

	// Extracting the refresh token and epoch time
	refreshToken.RefreshToken, refreshToken.RefreshTokenTime = fmt.Sprintf("%v", vaultData["refresh_token"]), fmt.Sprintf("%v", vaultData["refresh_epoch_time"])
	refereshTokenAgeLeft, err := strconv.ParseInt(refreshToken.RefreshTokenTime, 10, 64)
	if err != nil {
		log.Println("Failed to fetch refresh token")
		log.Println(err)
		validRefreshToken = false
		return "", validRefreshToken
	}
	if currTime > (refereshTokenAgeLeft + 8640000) {
		log.Println("Expired refresh token")
		log.Println(currTime, refereshTokenAgeLeft)
		validRefreshToken = false
		return "", validRefreshToken
	}

	return refreshToken.RefreshToken, validRefreshToken
}

func GetSNOWAuthToken() (SNOWResponse, error) {

	var snowResponse SNOWResponse
	refreshToken, validRefreshToken := GetSNOWRefreshToken()

	vaultData := vault.ReadSecrets("SNOW", "secret") // Fetching SNOW secrets

	if !validRefreshToken {

		payload := url.Values{}
		payload.Set("grant_type", "password")
		payload.Set("client_id", fmt.Sprintf("%v", vaultData["client_id"]))
		payload.Set("client_secret", fmt.Sprintf("%v", vaultData["client_secret"]))
		payload.Set("username", fmt.Sprintf("%v", vaultData["username"]))
		payload.Set("password", fmt.Sprintf("%v", vaultData["password"]))

		resp, err := isolatedfunctions.POSTforFormPayload(snowAuthEndpoint, payload) // Calls API with formdata and return http.Response and error
		if err != nil {
			log.Printf("POST call to dev instance failed due to %q", err)
			return snowResponse, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body) // Reads the body from http.Response and converts it into []byte
		// Parse JSON response

		if err := json.Unmarshal(body, &snowResponse); err != nil {
			log.Printf("Error decoding JSON: %q", err)
			return snowResponse, err
		} // Decodes string to json format

		vault.WriteSNOWRefreshToken(snowResponse.RefreshToken)

		return snowResponse, err
	}

	payload := url.Values{}
	payload.Set("grant_type", "refresh_token")
	payload.Set("client_id", fmt.Sprintf("%v", vaultData["client_id"]))
	payload.Set("client_secret", fmt.Sprintf("%v", vaultData["client_secret"]))
	payload.Set("refresh_token", refreshToken)

	resp, err := isolatedfunctions.POSTforFormPayload(snowAuthEndpoint, payload) // Calls API with formdata and return http.Response and error
	if err != nil {
		log.Printf("POST call to dev instance failed due to %q", err)
		return snowResponse, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body) // Reads the body from http.Response and converts it into []byte
	// Parse JSON response

	if err := json.Unmarshal(body, &snowResponse); err != nil {
		log.Printf("Error decoding JSON: %q", err)
		return snowResponse, err
	} // Decodes string to json format

	return snowResponse, err
}
