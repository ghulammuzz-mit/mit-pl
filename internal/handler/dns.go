package handler

import (
	"encoding/json"
	"io"
	"mit/platform/internal/cloudflare"
	"mit/platform/internal/dto"
	"mit/platform/internal/entity"
	"net/http"
)

func ListRecords(w http.ResponseWriter, r *http.Request) {
	cf := cloudflare.New()

	url := "https://api.cloudflare.com/client/v4/zones/" +
		cf.ZoneID + "/dns_records"

	req, _ := cf.NewRequest("GET", url, nil)
	resp, err := cf.HTTP.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Write(body)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	cf := cloudflare.New()

	var payload dto.CreateDNSRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, _ := json.Marshal(payload)
	url := "https://api.cloudflare.com/client/v4/zones/" +
		cf.ZoneID + "/dns_records"

	req, _ := cf.NewRequest("POST", url, body)
	resp, err := cf.HTTP.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	w.Write(respBody)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cf := cloudflare.New()

	body, _ := io.ReadAll(r.Body)
	url := "https://api.cloudflare.com/client/v4/zones/" +
		cf.ZoneID + "/dns_records/" + id

	req, _ := cf.NewRequest("PUT", url, body)
	resp, err := cf.HTTP.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	w.Write(respBody)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	recordName := r.PathValue("name")
	if recordName == "" {
		http.Error(w, "record name is required", http.StatusBadRequest)
		return
	}

	cf := cloudflare.New()

	searchURL := "https://api.cloudflare.com/client/v4/zones/" +
		cf.ZoneID + "/dns_records?name=" + recordName

	searchReq, _ := cf.NewRequest("GET", searchURL, nil)
	searchResp, err := cf.HTTP.Do(searchReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer searchResp.Body.Close()

	searchBody, _ := io.ReadAll(searchResp.Body)

	var listResp entity.CFListResponse
	if err := json.Unmarshal(searchBody, &listResp); err != nil {
		http.Error(w, "failed to parse cloudflare response", http.StatusInternalServerError)
		return
	}

	if len(listResp.Result) == 0 {
		http.Error(w, "dns record not found", http.StatusNotFound)
		return
	}

	recordID := listResp.Result[0].ID

	deleteURL := "https://api.cloudflare.com/client/v4/zones/" +
		cf.ZoneID + "/dns_records/" + recordID

	deleteReq, _ := cf.NewRequest("DELETE", deleteURL, nil)
	deleteResp, err := cf.HTTP.Do(deleteReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer deleteResp.Body.Close()

	deleteBody, _ := io.ReadAll(deleteResp.Body)
	w.WriteHeader(deleteResp.StatusCode)
	w.Write(deleteBody)
}
