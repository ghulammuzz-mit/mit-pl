package entity

type DNSRecord struct {
    ID      string `json:"id,omitempty"`
    Type    string `json:"type"`
    Name    string `json:"name"`
    Content string `json:"content"`
    TTL     int    `json:"ttl"`
    Proxied bool   `json:"proxied"`
}

type CFListResponse struct {
		Success bool `json:"success"`
		Result  []struct {
			ID string `json:"id"`
		} `json:"result"`
	}