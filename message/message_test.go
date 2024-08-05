package message

import (
	"bytes"
	"encoding/binary"
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
	sentMessage = []byte{0, 0, 0, 0}
	reader = bytes.NewReader(sentMessage)
	result, err = ReadMessage(reader)
	if err != nil {
		t.Errorf("Error reading result: %s", err)
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

func TestFormatHaveMessage(t *testing.T) {
	t.Log("Testing formatHaveMessage")
	expectedMessage := &Message{
		ID:      MsgHave,
		Payload: []byte{0x0, 0x0, 0x0, 0x2},
	}
	result := FormatHaveMessage(2)
	if !reflect.DeepEqual(result, expectedMessage) {
		t.Errorf("Expected have message: %v, got: %v", expectedMessage, result)
	}
}

func TestFormatRequest(t *testing.T) {
	t.Log("Testing Send Request")
	index := 2
	begin := 28
	length := 10
	expected := &Message{
		ID:      MsgRequest,
		Payload: make([]byte, 12),
	}
	binary.BigEndian.PutUint32(expected.Payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(expected.Payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(expected.Payload[8:12], uint32(length))
	result := FormatRequest(index, begin, length)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected: %v, got: %v", expected, result)

	}
}
