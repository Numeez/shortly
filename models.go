package shortly

import (
	"encoding/json"
	"log"
	"net/http"
)

type UrlRequest struct {
	Url string `json:"url"`
}

type URlShortenerResponse struct {
	Success bool               `json:"success"`
	Data    interface{}        `json:"data"`
	Error   *ShortenerAPIError `json:"error"`
}

type ShortenerAPIError struct {
	Message string `json:"message"`
}
type ShortenURL struct {
	Url string `json:"url"`
}

func ApiErrorMessage(err error, errorMessage string) *ShortenerAPIError {
	if err != nil {
		return &ShortenerAPIError{Message: err.Error()}
	}
	if errorMessage != "" {
		return &ShortenerAPIError{Message: errorMessage}
	}
	return nil
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string, err error) {
	apiError := ApiErrorMessage(err, errorMessage)
	response := URlShortenerResponse{
		Success: false,
		Data:    nil,
		Error:   apiError,
	}
	responseBytes, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseBytes)

}

func WriteResponse(w http.ResponseWriter, statusCode int, data any) {
	response := URlShortenerResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}
	responseBytes, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseBytes)
}
