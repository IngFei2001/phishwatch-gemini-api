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
	geminiKey := os.Getenv("GEMINI_API_KEY")
	fmt.Println("[PhishWatch] GEMINI_KEY length:", len(geminiKey))
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + geminiKey
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}
	jsonData, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("[PhishWatch] Gemini raw response:", string(respBody))
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("Gemini error: %s", string(respBody))
	}
	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content structure")
	}
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("invalid parts structure")
	}
	text, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text structure")
	}
	return text, nil
}
