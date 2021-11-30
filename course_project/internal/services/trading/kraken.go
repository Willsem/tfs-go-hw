package trading

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/dto"
)

const apiv3 = "/api/v3"

var (
	UnknownResponseError = errors.New("unknown response")
)

type KrakenTradingService struct {
	config config.Kraken
	client http.Client
}

func NewKrakenTradingService(config config.Kraken) *KrakenTradingService {
	return &KrakenTradingService{
		config: config,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

const (
	endpointOpenPositions = apiv3 + "/openpositions"
	methodOpenPositions   = http.MethodGet
)

func (service *KrakenTradingService) OpenPositions() ([]dto.Position, error) {
	res, err := service.sendRequest(methodOpenPositions, endpointOpenPositions, "")
	if err != nil {
		return nil, err
	}

	resp, ok := res["openPositions"]
	if !ok {
		return nil, UnknownResponseError
	}

	var positions []dto.Position
	j, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &positions)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

const (
	endpointSendOrder = apiv3 + "/sendorder"
	methodSendOrder   = http.MethodPost
)

func (service *KrakenTradingService) SendOrder(order dto.Order) (domain.OrderStatus, error) {
	postData := order.GetPostData()
	res, err := service.sendRequest(methodSendOrder, endpointSendOrder, postData)
	if err != nil {
		return domain.Empty, err
	}

	resp, ok := res["sendStatus"]
	if !ok {
		return domain.Empty, UnknownResponseError
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		return domain.Empty, UnknownResponseError
	}

	status, ok := respMap["status"]
	if !ok {
		return domain.Empty, UnknownResponseError
	}

	statusString, ok := status.(string)
	if !ok {
		return domain.Empty, UnknownResponseError
	}

	return domain.OrderStatus(statusString), nil
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

	if v, ok := result["result"]; !ok || v != "success" {
		errorMessage, ok := result["error"]
		if !ok {
			return nil, UnknownResponseError
		}

		return nil, errors.New(errorMessage.(string))
	}

	return result, nil
}
