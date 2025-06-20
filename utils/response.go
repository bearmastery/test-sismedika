package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON mengirimkan response HTTP dalam format JSON.
//
// Parameters:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - status: kode status HTTP (contoh: 200, 400, 500).
//   - data: objek data yang akan di-encode ke dalam format JSON.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
	}
}

// WriteError mengirimkan pesan error dalam format JSON standar.
//
// Parameters:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - status: kode status HTTP yang merepresentasikan jenis error.
//   - message: pesan error yang akan dikirim ke client dalam bentuk {"error": "..."}.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
