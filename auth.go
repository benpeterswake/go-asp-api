package asp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) refreshToken() {
	if c.authData.ExpiresIn == 0 {
		err := c.authenticate()
		if err != nil {
			log.Fatal(err)
		}
	}
	ticker := time.NewTicker(time.Duration(c.authData.ExpiresIn) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				log.Println("refreshing token")
				err := c.authenticate()
				ticker = time.NewTicker(time.Duration(c.authData.ExpiresIn) * time.Second)
				if err != nil {
					log.Fatal(err)
					close(quit)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *Client) authenticate() error {
	authEndpoint := "https://api.amazon.com/auth/o2/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.authData.RefreshToken)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	u, err := url.ParseRequestURI(authEndpoint)
	if err != nil {
		return err
	}
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	var response AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	c.authData = response
	return nil
}
