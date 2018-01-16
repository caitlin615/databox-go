package databox

import (
	_ "fmt"
	"os"
	"reflect"
	"testing"
	"time"
	_ "time"
)

var originalPostRequest = postRequest
var originalGetRequest = getRequest

func getToken() (pushToken string) {
	pushToken = "adxg1kq5a4g04k0wk0s4wkssow8osw84"
	envPushToken := "" + os.Getenv("DATABOX_PUSH_TOKEN")
	if envPushToken != "" {
		pushToken = os.Getenv("DATABOX_PUSH_TOKEN")
	}
	return
}

func TestSimpleInit(t *testing.T) {
	token := getToken()
	client := NewClient(getToken())

	if reflect.ValueOf(client).Kind().String() != "ptr" {
		t.Error("Not pointer")
	}

	if client.PushToken != token {
		t.Error("Token is not set.")
	}
}

func TestLastPush(t *testing.T) {
	lastPush, err := NewClient(getToken()).LastPush()
	if err != nil {
		t.Error("Error was raised", err)
	}

	if len(lastPush.Request.Errors) != 0 {
		t.Error("Number of errors in last push must equal 0!")
	}

	if len(lastPush.Request.Body.Data) == 0 {
		t.Error("Push must not be nil")
	}
}

func TestKPI_ToJsonData(t *testing.T) {
	a := (&KPI{Key: "a", Value: float32(33)}).ToJsonData()
	if a["$a"] != float32(33) {
		t.Error("Conversion error")
	}

	date := "2015-01-01 09:00:00"
	b := (&KPI{Key: "a", Date: date}).ToJsonData()
	if b["date"] != date {
		t.Error("Conversion error")
	}
}

func TestSuccessfulPush(t *testing.T) {
	if status, _ := NewClient(getToken()).Push(&KPI{
		Key:   "temp.ny",
		Value: 60.0,
	}); len(status.Errors) != 0 {
		t.Error("Not inserted")
	}
}

func TestFailedPush(t *testing.T) {
	if status, _ := NewClient(getToken()).Push(&KPI{
		Key:   "temp.ny",
		Value: 52.0,
		Date:  "2015-01-01 09:00:00",
	}); len(status.Errors) == 0 {
		// FIXME: This doesn't fail
		t.Error("This should not be \"ok\"")
	}
}

func TestWithAdditionalAttributes(t *testing.T) {
	postRequest = originalPostRequest
	getRequest = originalGetRequest

	client := NewClient(getToken())

	var attributes = make(map[string]interface{})
	attributes["test.number"] = 10
	attributes["test.string"] = "Oto Brglez"

	if status, _ := client.Push(&KPI{
		Key:        "test.TestWithAdditionalAttributes",
		Value:      10.0,
		Date:       time.Now().Format(time.RFC3339),
		Attributes: attributes,
	}); len(status.Errors) != 0 {
		t.Error("This status must be ok")
	}

	if _, err := client.LastPush(); err != nil {
		t.Error("Must be nil")
	}
}
