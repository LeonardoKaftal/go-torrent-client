package message

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadMessage(t *testing.T) {
	t.Log("Testing ReadMessage")
	sentMessage := []byte{0, 0, 0, 10, 4, 32, 43, 54, 65, 43, 23, 1, 44, 87}

	reader := bytes.NewReader(sentMessage)
	result, err := ReadMessage(reader)
	if err != nil {
		t.Errorf("Error reading result: %s", err)
	}
	expectedMessage := &Message{
		MsgHave,
		sentMessage[5:],
	}
	if !reflect.DeepEqual(result, expectedMessage) {
		t.Errorf("Expected: %v, got: %v", expectedMessage, result)
	}
}

func TestReadForAnInvalidMessage(t *testing.T) {
	t.Log("Testing ReadForAnInvalidMessage")
	sentMessage := []byte{0, 0, 0, 10, 4, 32, 43}
	reader := bytes.NewReader(sentMessage)
	_, err := ReadMessage(reader)
	if err == nil {
		t.Errorf("Expected an error for trying to read an invalid message (truncated), but got nil")
	}
}

func TestSerializeMessage(t *testing.T) {
	t.Log("Testing Serialize Message")
	messageToSent := &Message{
		ID:      MsgHave,
		Payload: []byte{32, 43, 54, 65, 43, 23, 1, 44, 87},
	}
	result := messageToSent.Serialize()
	expextedResult := []byte{0, 0, 0, 10, 4, 32, 43, 54, 65, 43, 23, 1, 44, 87}

	if !reflect.DeepEqual(result, expextedResult) {
		t.Errorf("Expected: %v, got: %v", expextedResult, result)
	}
}
