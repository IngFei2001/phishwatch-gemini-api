package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func askGemini(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	fmt.Println("[PhishWatch] API KEY length:", len(apiKey))
	url := "https://api.groq.com/openai/v1/chat/completions"

	body := map[string]interface{}{
		"model": "llama-3.1-8b-instant",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("[PhishWatch] Groq raw response:", string(respBody))

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("Groq error: %s", string(respBody))
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message structure")
	}

	text, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content structure")
	}

	return text, nil
}
