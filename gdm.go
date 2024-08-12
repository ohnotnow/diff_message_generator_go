package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type AIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type AIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json.maked_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func main() {
	model := flag.String("model", "gpt-4o-mini", "The model to use")
	flag.Parse()

	// Capture additional context if provided as arguments after the flags
	extraContext := strings.Join(flag.Args(), " ")

	cmd := exec.Command("git", "--no-pager", "diff", "--color-moved=no")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error running `git diff`:", stderr.String())
		return
	}

	systemMessage := `
You are an AI assistant who specializes in reading the output of 'git diff' and providing a well-written commit message to go with it.
Your commit message should cover *all* of the changes, not just the major ones.

Ensure you follow best practices for commit messages.

Guidelines:

- Capitalization and Punctuation: Capitalize the first word and do not end in punctuation. If using Conventional Commits, remember to use all lowercase.
- Mood: Use imperative mood in the subject line. Example: "Add fix for dark mode toggle state". Imperative mood gives the tone you are giving an order or request.
- Type of Commit: Specify the type of commit. It is recommended and can be even more beneficial to have a consistent set of words to describe your changes. Example: Bugfix, Update, Refactor, Bump, and so on.
- Length: The first line should ideally be no longer than 50 characters, and the body should be restricted to 72 characters.
- Content: Be direct, try to eliminate filler words and phrases in these sentences (examples: though, maybe, I think, kind of). Think like a journalist.
- Reference: If the user provides extra context including an issue number, include the issue number in the commit message. Example: "Fix #123: Add fix for dark mode toggle state".
- If there are a lot of changes use markdown bullet-points to list them. Example: "- Add fix for dark mode toggle state\n- Update the README.md file\n- Refactor the codebase\n- Bump the version to 1.0.0"
`

	prompt := fmt.Sprintf("I have the following output from running `git diff`. Could you give me a commit message for it? <diff>%s</diff>", stdout.String())
	if extraContext != "" {
		prompt = fmt.Sprintf("Context: %s\n\n%s", extraContext, prompt)
	}

	request := AIRequest{
		Model: *model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: prompt},
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to AI model:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var aiResponse AIResponse
	err = json.Unmarshal(responseBody, &aiResponse)
	if err != nil {
		fmt.Println("Error parsing AI response:", err)
		return
	}

	if len(aiResponse.Choices) == 0 {
		fmt.Println("No response from AI model")
		return
	}

	commitMessage := strings.TrimSpace(aiResponse.Choices[0].Message.Content)
	commitMessage = strings.Trim(commitMessage, "`\n ")
	fmt.Println(commitMessage)
}


