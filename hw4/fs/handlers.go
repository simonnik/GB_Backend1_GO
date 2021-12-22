package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const serveDir = "upload"
const hostAddr = "http://localhost:%d"

type File struct {
	Name string
	Ext  string
	Size int64
}

// upload handler that save files in the serveDir
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// get file from the form by field "file" from request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	// don't forget to close the file
	defer file.Close()

	// red file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	// check that the upload directory exists
	if _, err := os.Stat(serveDir); os.IsNotExist(err) {
		if err := os.Mkdir(serveDir, 0777); err != nil {
			log.Println(err)
			if err != nil {
				http.Error(w, "Cannot create dir", http.StatusInternalServerError)
				return
			}
		}
		log.Println(err)
	}

	filePath := serveDir + "/" + header.Filename

	// write file on server
	err = ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Unable to save file: %v", err), http.StatusInternalServerError)
		return
	}

	fileLink := fmt.Sprintf(hostAddr+"/"+header.Filename, fsPort)

	// send back a link of saved file
	fmt.Fprintln(w, fileLink)
}

// list handler that lists data about files stored in serveDir
func listHandler(w http.ResponseWriter, r *http.Request) {
	// check if method is not GET
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// get a list of filesystem objects on the server
	fileList, err := ioutil.ReadDir(serveDir)
	if err != nil {
		log.Printf("cannot read file from directory %s", serveDir)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
		return
	}

	// get extension filter from query
	filterExt := r.FormValue("ext")

	var files []File
	for _, file := range fileList {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())

			if filterExt == "" || ext == filterExt {
				fileAttr := File{
					Name: file.Name(),
					Ext:  ext,
					Size: file.Size(),
				}

				files = append(files, fileAttr)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		log.Printf("cannot convert data into json %v", err)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
		return
	}
}
