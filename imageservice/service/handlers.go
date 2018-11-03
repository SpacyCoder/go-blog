package service

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spacy/goblog/common/messaging"
)

var MessagingClient messaging.IMessagingClient

/**
 * Takes the POST body, decodes, processes and finally writes the result to the response.
 */
func ProcessImage(w http.ResponseWriter, r *http.Request) {

	sourceImage, _, err := image.Decode(r.Body)
	if err != nil {
		writeServerError(w, err.Error())
		return
	}
	writeAndReturn(w, sourceImage)
}

func GetAccountImage(w http.ResponseWriter, r *http.Request) {
	data := []byte("http://imageservice:7777/file/cake.jpg")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

/**
 * Takes the filename and tries to decode an image from /testimages/{filename}. Used for testing.
 */
func ProcessImageFromFile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var filename = vars["filename"]
	logrus.Println("Serving image for account: " + filename)

	fImg1, err := os.Open("testimages/" + filename)
	defer fImg1.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sourceImage, _, err := image.Decode(fImg1)

	if err != nil {
		writeServerError(w, err.Error())
		return
	}
	writeAndReturn(w, sourceImage)
}

func writeAndReturn(w http.ResponseWriter, sourceImage image.Image) {
	buf := new(bytes.Buffer)
	err := Sepia(sourceImage, buf)

	if err != nil {
		fmt.Println(err.Error())
		writeServerError(w, err.Error())
		return
	}
	outputData := buf.Bytes()

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(outputData)))
	w.WriteHeader(http.StatusOK)
	w.Write(outputData)

}

func writeServerError(w http.ResponseWriter, msg string) {
	logrus.Error(msg)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
