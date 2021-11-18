package trading

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
)

type KrakenTradingService struct {
	config config.Kraken
	client http.Client
}

func New(config config.Kraken) *KrakenTradingService {
	return &KrakenTradingService{
		config: config,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

const endpointOpenPositions = "/api/v3/openpositions"

func (service *KrakenTradingService) OpenPositions() error {
	res, err := service.sendRequest(http.MethodGet, endpointOpenPositions, "")
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

const endpointSendOrder = "/api/v3/sendorder"

func (service *KrakenTradingService) SendOrder(order Order) error {
	postData := ""
	_, err := service.sendRequest(http.MethodGet, endpointSendOrder, postData)
	if err != nil {
		return err
	}

	return nil
}

func (service *KrakenTradingService) CancelAllOrders() error {
	return nil
}

func (service *KrakenTradingService) generateAuthent(endpoint string, postData string) (string, error) {
	sha := sha256.New()
	sha.Write([]byte(postData + endpoint))

	decodedSecret, err := base64.StdEncoding.DecodeString(service.config.ApiSecret)
	if err != nil {
		return "", err
	}

	hmc := hmac.New(sha512.New, decodedSecret)
	hmc.Write(sha.Sum(nil))

	return base64.StdEncoding.EncodeToString(hmc.Sum(nil)), nil
}

func (service *KrakenTradingService) sendRequest(
	method string,
	endpoint string,
	postData string,
) (map[string]interface{}, error) {
	authent, err := service.generateAuthent(endpoint, postData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, service.config.RestUrl+endpoint+"?"+postData, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authent", authent)
	req.Header.Add("APIKey", service.config.ApiKey)

	resp, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
