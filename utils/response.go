package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// WriteJSON mengirimkan response HTTP dalam format JSON standar.
//
// Parameters:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - status: kode status HTTP (contoh: 200, 400, 500).
//   - data: objek data yang akan dikirim ke field "data".
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp := APIResponse{Data: data}
	buf, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(buf)
}

// WriteError mengirimkan pesan error dalam format JSON standar.
//
// Parameters:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - status: kode status HTTP yang merepresentasikan jenis error.
//   - message: pesan error yang akan dikirim ke client dalam field "error".
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")

	resp := APIResponse{Error: message}
	buf, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, `{"error": "failed to encode error response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(buf)
}
