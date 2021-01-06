package asp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/sts"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

var apiURL string = "https://sandbox.sellingpartnerapi-na.amazon.com"
var region string = "us-east-1"
var env string = os.Getenv("APP_ENV")

func init() {
	// if prod
	if env == "prod" || env == "dev" {
		apiURL = "https://sellingpartnerapi-na.amazon.com"
	}
}

type Client struct {
	apiURL       string
	authData     AuthResponse
	clientID     string
	clientSecret string
	roleArn      string
	roleID       string
}

// New creates a new Amazon Selling Partner Api Client
func New(clientID, clientSecret, refreshToken, roleArn, roleID string) *Client {
	return &Client{
		apiURL: apiURL,
		authData: AuthResponse{
			RefreshToken: refreshToken,
		},
		clientID:     clientID,
		clientSecret: clientSecret,
		roleArn:      roleArn,
		roleID:       roleID,
	}
}

func (c *Client) GetInventory() (*GetInventorySummariesResponse, error) {
	uri := "/fba/inventory/v1/summaries"
	// create request
	r, err := c.createRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	data := r.URL.Query()
	data.Set("granularityType", "Marketplace")
	data.Set("granularityId", "ATVPDKIKX0DER")
	data.Set("marketplaceIds", "ATVPDKIKX0DER")
	data.Set("details", "true")
	r.URL.RawQuery = data.Encode()

	r, err = c.signRequest(r)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return nil, errors.New(bodyString)
	}
	var response GetInventorySummariesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetOrders() error {
	uri := "/orders/v0/orders"
	r, err := c.createRequest(http.MethodGet, uri, nil)
	if err != nil {
		return err
	}
	data := r.URL.Query()
	data.Set("MarketplaceIds", "ATVPDKIKX0DER")
	r.URL.RawQuery = data.Encode()

	r, err = c.signRequest(r)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		return errors.New(bodyString)
	}

	return nil
}

func (c *Client) createRequest(method string, uri string, body io.Reader) (*http.Request, error) {
	c.refreshToken()
	ctx := context.Background()
	u, err := url.ParseRequestURI(c.apiURL + uri)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(method, u.String(), nil) // URL-encoded payload
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)
	// add headers
	r.Header.Add("content-type", "application/json")
	r.Header.Add("host", "https://sandbox.sellingpartnerapi-na.amazon.com")
	r.Header.Add("user-agent", "Sierra Home Systems Api/1.0 (Language=Golang/1.5;Platform=MacOS/11)")
	r.Header.Add("x-amz-access-token", c.authData.AccessToken)
	return r, nil
}

func (c *Client) signRequest(r *http.Request) (*http.Request, error) {
	ctx := context.Background()

	// Load config from env vars
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Assume role
	stsClient := sts.New(sts.Options{
		Region:      region,
		Credentials: cfg.Credentials,
	})

	params := &sts.AssumeRoleInput{
		RoleArn:         &c.roleArn,
		RoleSessionName: &c.roleID,
	}
	output, err := stsClient.AssumeRole(ctx, params)
	if err != nil {
		return nil, err
	}

	//get full credentials
	credOutput := credentials.NewStaticCredentialsProvider(*output.Credentials.AccessKeyId, *output.Credentials.SecretAccessKey, *output.Credentials.SessionToken)
	if err != nil {
		return nil, err
	}
	cred, err := credOutput.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	var bodyString string
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
	}
	// hash request body
	h := sha256.New()
	h.Write([]byte(bodyString))
	shaHash := hex.EncodeToString(h.Sum(nil))

	// create signer
	signer := v4.NewSigner()
	err = signer.SignHTTP(ctx, cred, r, shaHash, "execute-api", region, time.Now())
	if err != nil {
		fmt.Printf("failed to sign request: (%v)\n", err)
		return nil, err
	}
	return r, nil
}
