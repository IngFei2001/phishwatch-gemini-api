package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var GEMINI_KEY = os.Getenv("GEMINI_API_KEY")

func askGemini(prompt string) (string, error) {

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + GEMINI_KEY

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

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	candidates := result["candidates"].([]interface{})
	content := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	text := parts[0].(map[string]interface{})["text"].(string)

	return text, nil
}