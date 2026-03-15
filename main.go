package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AskRequest struct {
	URL       string   `json:"url"`
	RiskScore int      `json:"risk_score"`
	Signals   []string `json:"signals"`
	Question  string   `json:"question"`
}

type AskResponse struct {
	Answer string `json:"answer"`
}

func askAIHandler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("EXTENSION_TOKEN")
	fmt.Println("[PhishWatch] Request received")
	if r.Header.Get("X-EXTENSION-TOKEN") != token {
		fmt.Println("[PhishWatch] Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req AskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("[PhishWatch] Calling Gemini for URL:", req.URL)
	prompt := fmt.Sprintf(`
You are a Web3 phishing security assistant called PhishWatch AI.
Website: %s
Risk Score: %d
Signals: %s
User Question:
%s
Explain clearly whether the website is safe or risky.
Give a short security recommendation.
`, req.URL, req.RiskScore, strings.Join(req.Signals, ", "), req.Question)
	answer, err := askGemini(prompt)
	if err != nil {
		fmt.Println("[PhishWatch] Gemini error:", err)
		http.Error(w, "AI error", http.StatusInternalServerError)
		return
	}
	fmt.Println("[PhishWatch] Gemini success")
	resp := AskResponse{Answer: answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/ask-ai", askAIHandler)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
