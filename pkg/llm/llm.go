package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func formatPrompt(context string) string {
	return fmt.Sprintf(`
		As an English teacher, please correct and polish the following sentences while keeping their original meaning.  
		Return only the corrected version along with the original sentence.
		remove <think> in response
		Follow this format:

		1. **Original Sentence**: "Could yo hel me get this obj."  
		- **Revised**: "Can you help me get this object?"  

		2. **Original Sentence**: "I don't want to receive analysis instead of I want to."  
		- **Revised**: "I don't want to receive an analysis; instead, I want to."  

		3. **Original Sentence**: "Receiving corrected and polished version of the sentence."  
		- **Revised**: "A well-revised and polished version of this sentence is as follows."  

		Here is the sentence that needs correction:  
		"%s"
`, context)
}

func cleanResponse(response string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`) // Match anything inside <think>...</think>
	return re.ReplaceAllString(response, "") // Remove it
}

func CallLLM(context string) (string, error) {
	// Define request payload
	payload := map[string]interface{}{
		"model":  "deepseek-r1:1.5b",
		"prompt": formatPrompt(context),
		"stream": false,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Send HTTP POST request
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}

	return cleanResponse(response.Response), nil

}
