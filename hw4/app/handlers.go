package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/simonnik/GB_Backend1_GO/hw4/internal/models"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.FormValue("name")
		fmt.Fprintf(w, "Parsed query-param with key \"name\": %s\n", name)
	case http.MethodPost:
		var employee models.Employee

		contentType := r.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
		case "application/xml":
			err := xml.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal XML", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Got a new employee!\nName: %s\nAge: %dy.o.\nSalary %0.2f\n",
			employee.Name,
			employee.Age,
			employee.Salary,
		)
	}
}

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(h.UploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(h.UploadDir, 0777); err != nil {
			log.Println(err)
			if err != nil {
				http.Error(w, "Cannot create dir", http.StatusInternalServerError)
				return
			}
		}
		log.Println("stat error", err)
	}

	filepath := h.UploadDir + "/" + header.Filename

	err = ioutil.WriteFile(filepath, data, 0600)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Unable to save file: %v", err), http.StatusInternalServerError)
		return
	}

	fileLink := h.HostAddr + "/" + header.Filename
	log.Println("Send req: ", fileLink)

	req, err := http.NewRequest(http.MethodHead, fileLink, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to check file", http.StatusInternalServerError)
		return
	}

	cli := &http.Client{}

	resp, err := cli.Do(req)

	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to check file", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(resp)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, fileLink)
}
