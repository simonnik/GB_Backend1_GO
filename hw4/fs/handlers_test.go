package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadHandlerMethodNotAllowed(t *testing.T) {
	body := &bytes.Buffer{}
	req, _ := http.NewRequest(http.MethodGet, "/upload", body)
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestUploadHandlerUnableReadFile(t *testing.T) {
	body := &bytes.Buffer{}
	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	rr := httptest.NewRecorder()
	uploadHandler(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestUploadHandler(t *testing.T) {
	expected := `testfile.txt`
	file, _ := os.Open(expected)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	uploadHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestListHandler(t *testing.T) {
	body := &bytes.Buffer{}
	req, _ := http.NewRequest(http.MethodGet, "/list?ext=.txt", body)
	rr := httptest.NewRecorder()
	listHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `testfile.txt`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
