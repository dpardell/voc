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

// LLMClient is the interface satisfied by all LLM provider backends.
type LLMClient interface {
	GenerateQuiz(ctx context.Context, words []string, progress string, hostLang, targetLang string) ([]Question, error)
	Chat(ctx context.Context, progress string, history []Message, userMessage string, hostLang, targetLang string) (*ChatResponse, []Message, error)
	UpdateProgress(ctx context.Context, currentProgress string, sessionResults string, randomWords []string) (string, error)
	Close()
}

type GeminiClient struct {
	client    *genai.Client
	model     *genai.GenerativeModel
	modelName string
}

func NewGeminiClient(apiKey, projectID, location, modelName string) (*GeminiClient, error) {
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

	return &GeminiClient{
		client:    client,
		model:     model,
		modelName: modelName,
	}, nil
}

func (c *GeminiClient) Close() {
	c.client.Close()
}

func sanitizeJSON(text string) string {
	text = strings.TrimSpace(text)
	// Remove markdown code blocks if present
	if strings.Contains(text, "```") {
		parts := strings.SplitSeq(text, "```")
		for part := range parts {
			part = strings.TrimSpace(part)
			part, has := strings.CutPrefix(part, "json")
			if has {
				return strings.TrimSpace(part)
			}
			// If it's just a code block without "json" tag
			if len(part) > 0 && (strings.HasPrefix(part, "{") || strings.HasPrefix(part, "[")) {
				return part
			}
		}
	}
	return text
}

func (c *GeminiClient) generateWithRetry(ctx context.Context, model *genai.GenerativeModel, prompt string) (*genai.GenerateContentResponse, error) {
	var resp *genai.GenerateContentResponse
	var err error
	maxRetries := 3
	backoff := 2 * time.Second

	for i := range maxRetries {
		resp, err = model.GenerateContent(ctx, genai.Text(prompt))
		if err == nil {
			return resp, nil
		}

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

func (c *GeminiClient) GenerateQuiz(ctx context.Context, words []string, progress string, hostLang, targetLang string) ([]Question, error) {
	prompt := fmt.Sprintf(`
You are a language teacher designing intentional grammar and vocabulary practice.

Host Language (Learner's native/UI language): %s
Target Language (Language being learned): %s
Vocabulary pool (use these words as raw material): %s

Learner's current progress:
%s

Generate 5 multiple-choice questions that test grammar patterns and contextual usage, NOT simple word recall or translation. Use the vocabulary words as vehicles for testing the grammar concepts marked as Struggling and Progressing above.

Question design principles:
- Prefer sentence-completion questions: write a sentence in %s with a blank (use "___"), and provide 4 choices that differ in grammar (e.g. conjugation, tense, agreement, preposition). The learner picks the correct form.
- Also include questions that test usage in context: show a sentence using one of the words and ask something about the grammar (e.g. "What tense is used here?", "Which option correctly completes this sentence?").
- Avoid simple "what does X mean?" translation questions — these test memory, not understanding.
- Make distractors plausible: wrong answers should be grammatically close (wrong tense, wrong gender agreement, wrong preposition), not random.
- Write question text in %s when asking about grammar concepts; write it in %s when testing comprehension in context.

All 5 questions must be "multiple_choice" with exactly 4 options. Each must include:
- "type": "multiple_choice"
- "question": the question or sentence with blank
- "options": array of 4 strings
- "correct_answer_index": 0-based index of the correct option
- "correct_answer": the literal text of the correct option
- "target_word": the word from the vocabulary pool this question relates to

Return ONLY a valid JSON array of these 5 question objects. No preamble or markdown.
`, hostLang, targetLang, strings.Join(words, ", "), progress, targetLang, hostLang, targetLang)

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

func (c *GeminiClient) UpdateProgress(ctx context.Context, currentProgress string, sessionResults string, randomWords []string) (string, error) {
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
Produce a progress file with exactly this structure:

## Summary
2-3 sentences: estimated CEFR level, overall trajectory, notable strengths and gaps.

## Mastered
- **<Concept>** — <one-line note>
  - ✓ <Specific case>

## Progressing
- **<Concept>** — <one-line note on where they slip up>
  - ~ <Specific case they are working on>

## Struggling
- **<Concept>** — <one-line note on the core confusion>
  - ✗ <Specific case that keeps going wrong>

Rules:
- Concepts are grammar topics or patterns (e.g. "plus-que-parfait formation", "gérondif vs. participe présent").
- Specific cases are narrower instances within a concept (e.g. "si clause: plus-que-parfait + conditionnel passé").
- Each concept must have at least one specific case nested under it.
- Place every concept in exactly one of the three tiers based on the session results and existing progress.
- If existing progress uses a different format, migrate it into this structure.
- Do not list individual words — describe grammar concepts and patterns only.
- Only return the file content. No preamble or explanation.

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

func (c *GeminiClient) Chat(ctx context.Context, progress string, history []Message, userMessage string, hostLang, targetLang string) (*ChatResponse, []Message, error) {
	chatModel := c.client.GenerativeModel(c.modelName)
	chatModel.ResponseMIMEType = "application/json"

	chatModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(fmt.Sprintf(`
You are a professional language coach named Coach.

Learner Progress (structured into Mastered / Progressing / Struggling concepts):
%s

Use this structure to tailor conversation difficulty and focus corrections on Struggling and Progressing concepts.

TASK:
Converse with the user in %s.
Provide helpful grammar and spelling corrections for ONLY the user's very last message in the "corrections" field.
Explain any major mistakes or corrections briefly in %s.

CRITICAL RULES:
- You MUST return a JSON object. 
- Do NOT return plain text. 
- The "response" field must contain your reply in %s.
- The "corrections" field MUST be an array of objects, each with "incorrect" and "correct" keys.
  Example: [{"incorrect": "Slaut", "correct": "Salut"}]
- If there are no corrections needed, leave the "corrections" array empty: [].
`, progress, targetLang, hostLang, targetLang))},
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

	for i := range maxRetries {
		resp, err = session.SendMessage(ctx, genai.Text(userMessage))
		if err == nil {
			break
		}

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
