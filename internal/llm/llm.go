package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type Question struct {
	Type               string   `json:"type"` // "multiple_choice" or "fill_in_the_blank"
	Question           string   `json:"question"`
	Options            []string `json:"options,omitempty"`
	CorrectAnswerIndex int      `json:"correct_answer_index,omitempty"`
	CorrectAnswer      string   `json:"correct_answer"`
	TargetWord         string   `json:"target_word"`
}

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
	modelName string
}

func NewClient(apiKey, projectID, location, modelName string) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	if modelName == "" {
		modelName = "gemini-2.5-flash-lite"
	}

	model := client.GenerativeModel(modelName)
	model.ResponseMIMEType = "application/json"

	return &Client{
		client:    client,
		model:     model,
		modelName: modelName,
	}, nil
}


func (c *Client) Close() {
	c.client.Close()
}

func sanitizeJSON(text string) string {
	text = strings.TrimSpace(text)
	// Remove markdown code blocks if present
	if strings.Contains(text, "```") {
		parts := strings.Split(text, "```")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "json") {
				return strings.TrimSpace(strings.TrimPrefix(part, "json"))
			}
			// If it's just a code block without "json" tag
			if len(part) > 0 && (strings.HasPrefix(part, "{") || strings.HasPrefix(part, "[")) {
				return part
			}
		}
	}
	return text
}

func (c *Client) generateWithRetry(ctx context.Context, model *genai.GenerativeModel, prompt string) (*genai.GenerateContentResponse, error) {
	var resp *genai.GenerateContentResponse
	var err error
	maxRetries := 3
	backoff := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err = model.GenerateContent(ctx, genai.Text(prompt))
		if err == nil {
			return resp, nil
		}

		// Check if it's a 429 error
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusTooManyRequests {
			if i < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
		}
		return nil, err
	}
	return nil, err
}

func (c *Client) GenerateQuiz(ctx context.Context, words []string, progress string) ([]Question, error) {
	prompt := fmt.Sprintf(`
You are a language teacher creating an interactive vocabulary quiz.
Generate a quiz for the following target words: %s.

Current learner progress:
%s

Instructions:
1. Include two types of questions: "multiple_choice" and "fill_in_the_blank".
2. Ensure a mix of questions. Some should be translation-based (English to target language or vice versa), 
   while others should be entirely in the target language (e.g., fill in the blank in a sentence, 
   choosing the correct conjugation, or matching a target-language definition to a word).
3. "multiple_choice" questions should have:
   - "question": The question text.
   - "options": An array of 4 possible answers.
   - "correct_answer_index": The 0-based index of the correct answer.
   - "correct_answer": The literal text of the correct answer.
4. "fill_in_the_blank" questions should have:
   - "question": A sentence with a blank (use "___") or a prompt.
   - "correct_answer": The exact word or phrase that fills the blank.
   - "options": Should be null or empty.
5. All questions must include:
   - "type": "multiple_choice" or "fill_in_the_blank".
   - "target_word": The word from the target list that this question is testing.

Focus on variety and challenge, using the learner's progress to inform the difficulty level and which words to prioritize.
Return ONLY a valid JSON array of these question objects. No preamble or markdown formatting.
`, strings.Join(words, ", "), progress)

	resp, err := c.generateWithRetry(ctx, c.model, prompt)
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

	sanitized := sanitizeJSON(string(text))
	var questions []Question
	if err := json.Unmarshal([]byte(sanitized), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse quiz JSON: %v", err)
	}

	return questions, nil
}

type Correction struct {
	Incorrect string `json:"incorrect"`
	Correct   string `json:"correct"`
}

type ChatResponse struct {
	Response    string       `json:"response"`
	Corrections []Correction `json:"corrections"` // Suggestions or corrections for the user's last message
}

func (c *Client) UpdateProgress(ctx context.Context, currentProgress string, sessionResults string, randomWords []string) (string, error) {
	// Create a temporary model without JSON constraint for the progress update
	updateModel := c.client.GenerativeModel(c.modelName)
	updateModel.ResponseMIMEType = "text/plain"

	prompt := fmt.Sprintf(`
You are an expert language learning assistant tracking a user's progress.
Your task is to update the learner's progress file (in Markdown format) based on their latest session results.

Current Progress:
%s

Latest Session Results:
%s

Random sample of user's words for context:
%s

Instructions:
1. Update the Markdown progress file to provide a HIGH-LEVEL overview of the learner's skills and standing.
2. Focus on grammar competency, vocabulary range, estimated CEFR level, and specific strengths/weaknesses.
3. Don't just list words; describe the *types* of words and concepts they are mastering.
4. Be concise but insightful.
5. Maintain a structured format (e.g., using headers, lists).
6. Only return the NEW content of the progress file. Do not include any preamble or explanation.

Updated Progress File:
`, currentProgress, sessionResults, strings.Join(randomWords, ", "))

	resp, err := c.generateWithRetry(ctx, updateModel, prompt)
	if err != nil {
		return "", err
	}
	// ... rest of the function (omitted for brevity, assume it's correctly handled by the tool)

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return string(text), nil
}

type Message struct {
	Role    string // "user" or "model"
	Content string
}

func (c *Client) GenerateDailySaying(ctx context.Context, lang string, progress string) (string, error) {
	// Use plain text for this
	model := c.client.GenerativeModel(c.modelName)
	model.ResponseMIMEType = "text/plain"

	prompt := fmt.Sprintf(`
You are a language learning motivator.
The learner is practicing: %s.
Their current progress is:
%s

Generate a fun, short, and inspiring saying or idiom in the target language (%s) that is appropriate for their level.
Provide the saying, its pronunciation (if helpful), and a brief English translation.
Keep it under 3 lines.
`, lang, progress, lang)

	resp, err := c.generateWithRetry(ctx, model, prompt)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return string(text), nil
}

func (c *Client) Chat(ctx context.Context, progress string, history []Message, userMessage string) (*ChatResponse, []Message, error) {
	// Use gemini-2.0-flash for potentially better stability and speed
	chatModel := c.client.GenerativeModel(c.modelName)
	chatModel.ResponseMIMEType = "application/json"

	chatModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(fmt.Sprintf(`
You are a professional language coach.
Learner Progress: %s

TASK:
Converse with the user in their target language. 
Provide helpful grammar and spelling corrections for ONLY the user's very last message in the "corrections" field. 

CRITICAL RULES:
- You MUST return a JSON object. 
- Do NOT return plain text. 
- The "response" field must contain your reply in the target language.
- The "corrections" field MUST be an array of objects, each with "incorrect" and "correct" keys.
  Example: [{"incorrect": "Slaut", "correct": "Salut"}]
- If there are no corrections needed, leave the "corrections" array empty: [].
`, progress))},
	}

	session := chatModel.StartChat()

	// Convert history to genai.Content
	var genaiHistory []*genai.Content
	for _, msg := range history {
		role := msg.Role
		if role == "" {
			role = "user"
		}
		genaiHistory = append(genaiHistory, &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		})
	}
	session.History = genaiHistory

	var resp *genai.GenerateContentResponse
	var err error
	maxRetries := 3
	backoff := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err = session.SendMessage(ctx, genai.Text(userMessage))
		if err == nil {
			break
		}

		// Check if it's a 429 error
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusTooManyRequests {
			if i < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
		}
		return nil, nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, nil, fmt.Errorf("no response from LLM")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response format")
	}

	sanitized := sanitizeJSON(string(text))
	var chatRes ChatResponse
	if err := json.Unmarshal([]byte(sanitized), &chatRes); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON response: %v\nResponse was: %s", err, sanitized)
	}

	// Update history
	newHistory := append(history, Message{Role: "user", Content: userMessage})
	newHistory = append(newHistory, Message{Role: "model", Content: chatRes.Response})

	return &chatRes, newHistory, nil
}
