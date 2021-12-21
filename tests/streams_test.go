package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/asumsi/api.inlive/internal/models/stream"
)

func createStream(t *testing.T, streamModel *stream.Stream) stream.StreamResponse {

	bodyJson, err := json.Marshal(streamModel)
	body := bytes.NewBuffer(bodyJson)

	if err != nil {
		t.Errorf("Error on encode JSON")
	}
	req, _ := http.NewRequest("POST", "/v1/streams/create", body)

	resp := ExecuteRequest(App.router, req)
	ExpectResponseCode(t, http.StatusOK, resp.Code)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Failed on read body, got: %s", respBody)
	}

	var respBodyJson stream.StreamResponse
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Response not a JSON object, got: %s", respBody)
	}
	return respBodyJson
}
func TestStreamsCreate(t *testing.T) {
	streamModel := &stream.Stream{
		Name: "Test stream",
		Slug: "test-slug",
	}
	respBodyJson := createStream(t, streamModel)

	//need to test if the name from response is the same from stream that we post
	if respBodyJson.Data.Name != streamModel.Name {
		t.Errorf("Expect resp field value %s, got: %s", streamModel.Name, respBodyJson.Data.Name)
	}
}

func TestStreamsUpdate(t *testing.T) {
	streamModel := &stream.Stream{
		Name: "Test stream",
		Slug: "test-slug",
	}

	updatedName := "Name updated"
	updatedSlug := "test-slug-updated"

	createdJson := createStream(t, streamModel)
	updatedStream := createdJson.Data

	updatedStream.Name = updatedName
	updatedStream.Slug = updatedSlug

	bodyJson, err := json.Marshal(updatedStream)
	body := bytes.NewBuffer(bodyJson)

	if err != nil {
		t.Errorf("Error on encode JSON")
	}
	routepath := "/v1/streams/" + strconv.FormatInt(createdJson.Data.ID, 10)
	req, _ := http.NewRequest("PUT", routepath, body)

	resp := ExecuteRequest(App.router, req)
	ExpectResponseCode(t, http.StatusOK, resp.Code)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Failed on read body, got: %s", respBody)
	}

	var respBodyJson stream.StreamResponse
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Response not a JSON object, got: %s", respBody)
	}
	if respBodyJson.Data.Name != updatedName {
		t.Errorf("Expect resp field value %s, got: %s", updatedName, respBodyJson.Data.Name)
	}

	if respBodyJson.Data.Slug != updatedSlug {
		t.Errorf("Expect resp field value %s, got: %s", updatedSlug, respBodyJson.Data.Slug)
	}
}

func TestStreamsGet(t *testing.T) {
	streamModel := &stream.Stream{
		Name: "Test stream-get",
		Slug: "test-slug-get",
	}

	createdJson := createStream(t, streamModel)

	routepath := "/v1/streams/" + strconv.FormatInt(createdJson.Data.ID, 10)

	req, _ := http.NewRequest("GET", routepath, nil)

	resp := ExecuteRequest(App.router, req)
	ExpectResponseCode(t, http.StatusOK, resp.Code)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Failed on read body, got: %s", respBody)
	}

	var respBodyJson stream.StreamResponse
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Response not a JSON object, got: %s", respBody)
	}
	if respBodyJson.Data.ID != createdJson.Data.ID {
		t.Errorf("Expect resp field value %s, got: %s", strconv.FormatInt(respBodyJson.Data.ID, 10), strconv.FormatInt(createdJson.Data.ID, 10))
	}
}

func TestStreamsDelete(t *testing.T) {
	streamModel := &stream.Stream{
		Name: "Test stream-delete",
		Slug: "test-slug-delete",
	}

	createdJson := createStream(t, streamModel)

	routepath := "/v1/streams/" + strconv.FormatInt(createdJson.Data.ID, 10)

	req, _ := http.NewRequest("DELETE", routepath, nil)

	resp := ExecuteRequest(App.router, req)
	ExpectResponseCode(t, http.StatusOK, resp.Code)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Failed on read body, got: %s", respBody)
	}

	var respBodyJson stream.StreamResponse
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Response not a JSON object, got: %s", respBody)
	}
	if respBodyJson.Data.ID != createdJson.Data.ID {
		t.Errorf("Expect resp field value %s, got: %s", strconv.FormatInt(respBodyJson.Data.ID, 10), strconv.FormatInt(createdJson.Data.ID, 10))
	}
}

func TestStreamsInit(t *testing.T) {
	streamModel := &stream.Stream{
		Name: "Test stream-delete",
		Slug: "test-slug-delete",
	}

	createdJson := createStream(t, streamModel)

	routepath := "/v1/streams/" + strconv.FormatInt(createdJson.Data.ID, 10) + "/init"

	req, _ := http.NewRequest("POST", routepath, nil)

	resp := ExecuteRequest(App.router, req)
	ExpectResponseCode(t, http.StatusOK, resp.Code)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Failed on read body, got: %s", respBody)
	}

	var respBodyJson stream.StreamResponse
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil { // Parse []byte to the go struct pointer
		t.Errorf("Response not a JSON object, got: %s", respBody)
	}
	// if respBodyJson.Data.ID != webrtc.SessionDescription {
	// 	t.Errorf("Expect resp field value %s, got: %s", strconv.FormatInt(respBodyJson.Data.ID, 10), strconv.FormatInt(createdJson.Data.ID, 10))
	// }
}
