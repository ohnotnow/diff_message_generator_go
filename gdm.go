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

        // Get the output of git diff
    diffOutput := stdout.String()

    // Check if the git diff output is empty
    if diffOutput == "" {
        fmt.Println("No changes detected. There is nothing to commit.")
        return
    }

	systemMessage := `
You are an AI assistant specialized in reading the output of 'git diff' and generating well-structured commit messages following the **Conventional Commits** specification.

Your commit message should include:

1. **Subject Line**: Adhere to the Conventional Commits format.
2. **Summary**: Provide a concise summary of what the commit aims to achieve. If the user supplies additional context or guidance, incorporate that information; otherwise, make an educated guess based on the 'git diff' output.
3. **Detailed Changes**: List the main changes in a markdown bullet-point format.

**Guidelines:**

- **Subject Line (Conventional Commits)**:
  - Format: '<type>(optional scope): <description>'
  - **Type**: Use a consistent set of commit types such as 'feat', 'fix', 'docs', 'style', 'refactor', 'test', 'chore', etc.
  - **Scope**: (Optional) Specify the scope of the changes, e.g., 'auth', 'UI', 'database'.
  - **Description**: Use the imperative mood and all lowercase. Do not end with punctuation.
  - **Example**: 'feat(auth): add OAuth2 login functionality'

- **Summary**:
  - Provide a brief overview of the commit’s purpose.
  - If user-provided context is available, incorporate it to enhance accuracy.
  - Aim for clarity and conciseness.

- **Detailed Changes**:
  - Use markdown bullet points to enumerate the main changes.
  - Ensure each bullet point starts with a verb in the imperative mood.
  - Example:
    - Add OAuth2 login functionality
    - Update the authentication middleware
    - Refactor user session management
    - Improve error handling for login failures

**Additional Guidelines:**

- **Response**:
  - Only respond with the commit message, no other text or chat.
  
- **Capitalization and Punctuation**:
  - Capitalize the first word of the subject line.
  - Do not end the subject line with punctuation.
  
- **Length**:
  - Subject line: Ideally no longer than 50 characters.
  - Summary: Keep it concise, typically one to two sentences.
  - Detailed changes: Each bullet point should be clear and succinct.

- **Content**:
  - Be direct and eliminate filler words and phrases.
  - Think like a journalist—focus on the "who, what, why" of the changes.

- **References**:
  - If the user provides extra context, such as an issue number, include it in the subject line.
  - Example: 'fix(auth): resolve login bug causing session timeout (#123)'

**Example Commit Message:**

feat(auth): add OAuth2 login functionality

Introduce OAuth2 authentication to enhance security and provide third-party login options.

- Implement OAuth2 login endpoints
- Update authentication middleware to handle OAuth2 tokens
- Refactor user session management for OAuth2 compatibility
- Improve error handling for OAuth2 login failures
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


