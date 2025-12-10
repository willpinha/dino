package dino

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
	Age   int    `json:"age" xml:"age"`
}

func TestReadJSON_Success(t *testing.T) {
	jsonData := `{"name":"John","email":"john@example.com","age":30}`
	reader := strings.NewReader(jsonData)

	result, err := ReadJSON[testStruct](reader)

	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{"malformed JSON", `{"name":"John","email":}`},
		{"invalid syntax", `{name:"John"}`},
		{"incomplete", `{"name":"John"`},
		{"not JSON", `this is not JSON`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.jsonData)

			result, err := ReadJSON[testStruct](reader)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid JSON body")
			assert.Empty(t, result.Name)
		})
	}
}

func TestReadXML_Success(t *testing.T) {
	xmlData := `<testStruct><name>John</name><email>john@example.com</email><age>30</age></testStruct>`
	reader := strings.NewReader(xmlData)

	result, err := ReadXML[testStruct](reader)

	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)
}

func TestReadXML_InvalidXML(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
	}{
		{"invalid syntax", `<testStruct><name>John</testStruct>`},
		{"not XML", `this is not XML`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.xmlData)

			result, err := ReadXML[testStruct](reader)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid XML body")
			assert.Empty(t, result.Name)
		})
	}
}

func TestReadXML_EmptyBody(t *testing.T) {
	reader := strings.NewReader("")

	result, err := ReadXML[testStruct](reader)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid XML body")
	assert.Empty(t, result.Name)
}

func TestReadBytes_Success(t *testing.T) {
	data := []byte("Hello, World!")
	reader := bytes.NewReader(data)

	result, err := ReadBytes(reader)

	require.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestReadBytes_EmptyBody(t *testing.T) {
	reader := bytes.NewReader([]byte{})

	result, err := ReadBytes(reader)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestReadJSON_FromRequest(t *testing.T) {
	jsonData := `{"name":"John","email":"john@example.com","age":30}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(jsonData))

	result, err := ReadJSON[testStruct](req.Body)

	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)
}

func TestReadXML_FromRequest(t *testing.T) {
	xmlData := `<testStruct><name>John</name><email>john@example.com</email><age>30</age></testStruct>`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(xmlData))

	result, err := ReadXML[testStruct](req.Body)

	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)
}

func TestReadBytes_FromRequest(t *testing.T) {
	data := []byte("test data")
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(data))

	result, err := ReadBytes(req.Body)

	require.NoError(t, err)
	assert.Equal(t, data, result)
}

// errorReader is a helper type that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestReadJSON_ErrorDetails(t *testing.T) {
	jsonData := `invalid json`
	reader := strings.NewReader(jsonData)

	_, err := ReadJSON[testStruct](reader)

	require.Error(t, err)
	httpErr, ok := err.(*Error)
	require.True(t, ok)
	assert.NotNil(t, httpErr.Details)
}

func TestReadXML_ErrorDetails(t *testing.T) {
	xmlData := `invalid xml`
	reader := strings.NewReader(xmlData)

	_, err := ReadXML[testStruct](reader)

	require.Error(t, err)
	httpErr, ok := err.(*Error)
	require.True(t, ok)
	assert.NotNil(t, httpErr.Details)
}

func TestReadBytes_ErrorDetails(t *testing.T) {
	reader := &errorReader{}

	_, err := ReadBytes(reader)

	require.Error(t, err)
	httpErr, ok := err.(*Error)
	require.True(t, ok)
	assert.NotNil(t, httpErr.Details)
}

func TestReadJSON_PointerTypes(t *testing.T) {
	jsonData := `{"name":"John","email":"john@example.com","age":30}`
	reader := strings.NewReader(jsonData)

	result, err := ReadJSON[*testStruct](reader)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)
}

func TestReadJSON_WithExtraFields(t *testing.T) {
	jsonData := `{"name":"John","email":"john@example.com","age":30,"extra":"ignored"}`
	reader := strings.NewReader(jsonData)

	result, err := ReadJSON[testStruct](reader)

	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)
}
