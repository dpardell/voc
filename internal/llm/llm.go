package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Question struct {
	Question           string   `json:"question"`
	Options            []string `json:"options"`
	CorrectAnswerIndex int      `json:"correct_answer_index"`
}

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"

	return &Client{
		client: client,
		model:  model,
	}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) GenerateQuiz(words []string, contextWords []string) ([]Question, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf(`
You are a language teacher creating a multiple-choice vocabulary quiz.
Generate a quiz for the following target words: %s.
Use the following words as additional context or distractors if helpful: %s.

Return a JSON array of objects, where each object has:
- "question": The question text (e.g., "What is the meaning of 'Word'?" or "Which word means 'Definition'?").
- "options": An array of 4 possible answers.
- "correct_answer_index": The 0-based index of the correct answer in the options array.

Ensure the questions are suitable for someone learning the language.
Focus on the definitions and usage of the target words.
`, strings.Join(words, ", "), strings.Join(contextWords, ", "))

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from LLM")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var questions []Question
	if err := json.Unmarshal([]byte(text), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return questions, nil
}
