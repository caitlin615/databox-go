package databox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

const (
	apiURL        = "https://push2new.databox.com"
	clientVersion = "0.1.2"
)

// Client is the client used to communicate with the databox API
type Client struct {
	PushToken string
	PushHost  string
}

// KPI is a key performance indicator that is sent to the databox API
type KPI struct {
	Key        string
	Value      float32
	Date       string
	Attributes map[string]interface{}
}

// KPIWrap is used to serialize KPIs when se
type KPIWrap struct {
	Data []map[string]interface{} `json:"data"`
}

// NewClient creates a new Client with the supplied push token
func NewClient(pushToken string) *Client {
	return &Client{
		PushToken: pushToken,
		PushHost:  apiURL,
	}
}

// ResponseStatus is the response returned from the push API
type ResponseStatus struct {
	ID      string   `json:"id"`
	Metrics []string `json:"metrics"`
	Errors  []string `json:"errors"`
}

var postRequest = func(client *Client, path string, payload []byte) ([]byte, error) {
	userAgent := "Databox/" + clientVersion + " (" + runtime.Version() + ")"
	request, err := http.NewRequest("POST", (client.PushHost + path), bytes.NewBuffer(payload))
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/vnd.databox.v2+json")
	request.SetBasicAuth(client.PushToken, "")

	if err != nil {
		return nil, err
	}

	response, err2 := (&http.Client{}).Do(request)
	if err2 != nil {
		return nil, err2
	}

	data, err3 := ioutil.ReadAll(response.Body)
	return data, err3
}

var getRequest = func(client *Client, path string) ([]byte, error) {
	userAgent := "Databox/" + clientVersion + " (" + runtime.Version() + ")"
	request, err := http.NewRequest("GET", (client.PushHost + path), nil)
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/vnd.databox.v2+json")
	request.SetBasicAuth(client.PushToken, "")

	if err != nil {
		return nil, err
	}

	response, err2 := (&http.Client{}).Do(request)
	if err2 != nil {
		return nil, err2
	}

	data, err3 := ioutil.ReadAll(response.Body)
	return data, err3
}

// LastPush contains the data returned from the /lastpushes endpoint
type LastPush struct {
	Request struct {
		Date   time.Time `json:"date"`
		Errors []string  `json:"errors"` // FIXME: When no errors, returns "[]", but what happens if there are errors?
		Body   struct {
			Data []map[string]interface{} `json:"data"`
		} `json:"body"`
	} `json:"request"`
	Response struct {
		Date time.Time `json:"date"`
		Body struct {
			ID string `json:"string"`
		} `json:"body"`
	} `json:"response"`
	Metrics []string `json:"metrics"`
}

/*
TODO: Use this for better and fine grained response parsing
func (lastPush *LastPush) UnmarshalJSON(raw []byte) error {
	var preParsed map[string]interface{}
	if err := json.Unmarshal(raw, &preParsed); err != nil {
		return err
	}

	parsedTime, err1 := time.Parse(time.RFC3339, preParsed["datetime"].(string))
	if err1 != nil {
		return err1
	}

	lastPush.Time = parsedTime
	return nil
}
*/

// LastPushes returns info on the last n pushes
func (client *Client) LastPushes(n int) ([]LastPush, error) {
	response, err := getRequest(client, fmt.Sprintf("/lastpushes/%d", n))
	if err != nil {
		return nil, err
	}

	lastPushes := make([]LastPush, 0)
	err1 := json.Unmarshal(response, &lastPushes)
	if err1 != nil {
		return nil, err1
	}

	return lastPushes, nil
}

// LastPush returns the result from the /lastpushes endpoint
func (client *Client) LastPush() (LastPush, error) {
	lastPush := LastPush{}
	response, err := getRequest(client, "/lastpushes")
	if err != nil {
		return lastPush, err
	}

	lastPushes := make([]LastPush, 0)
	err1 := json.Unmarshal(response, &lastPushes)
	if err1 != nil {
		return lastPush, err1
	}
	if len(lastPushes) == 0 {
		return lastPush, fmt.Errorf("no pushes recorded")
	}

	return lastPushes[0], nil
}

// Push sends data to the databox API
func (client *Client) Push(kpi *KPI) (*ResponseStatus, error) {
	payload, err := serializeKPIs([]KPI{*kpi})
	if err != nil {
		fmt.Println("serialize")
		return &ResponseStatus{}, err
	}

	response, err2 := postRequest(client, "/", payload)
	if err2 != nil {
		fmt.Println("post response:")
		return &ResponseStatus{}, err2
	}

	var responseStatus = &ResponseStatus{}
	if err3 := json.Unmarshal(response, &responseStatus); err3 != nil {
		return &ResponseStatus{}, err3
	}

	return responseStatus, nil
}

// ToJSONData serializes KPI into a payload meant for the databox API
func (kpi *KPI) ToJSONData() map[string]interface{} {
	var payload = make(map[string]interface{})
	payload["$"+kpi.Key] = kpi.Value

	if kpi.Date != "" {
		payload["date"] = kpi.Date
	}

	if len(kpi.Attributes) != 0 {
		for key, value := range kpi.Attributes {
			payload[key] = value
		}
	}

	return payload
}

func serializeKPIs(kpis []KPI) ([]byte, error) {
	wrap := KPIWrap{
		Data: make([]map[string]interface{}, 0),
	}

	for _, kpi := range kpis {
		wrap.Data = append(wrap.Data, kpi.ToJSONData())
	}

	return json.Marshal(wrap)
}
