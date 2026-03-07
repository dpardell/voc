package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const mistralAPIBase = "https://api.mistral.ai/v1/chat/completions"

type MistralClient struct {
	apiKey     string
	modelName  string
	httpClient *http.Client
}

func NewMistralClient(apiKey, modelName string) *MistralClient {
	if modelName == "" {
		modelName = "mistral-small-latest"
	}
	return &MistralClient{
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{},
	}
}

func (c *MistralClient) Close() {}

type mistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type mistralRequest struct {
	Model          string           `json:"model"`
	Messages       []mistralMessage `json:"messages"`
	ResponseFormat *mistralFormat   `json:"response_format,omitempty"`
}

type mistralFormat struct {
	Type string `json:"type"`
}

type mistralResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *MistralClient) complete(ctx context.Context, messages []mistralMessage, jsonMode bool) (string, error) {
	req := mistralRequest{
		Model:    c.modelName,
		Messages: messages,
	}
	if jsonMode {
		req.ResponseFormat = &mistralFormat{Type: "json_object"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	maxRetries := 3
	backoff := 2 * time.Second

	for i := range maxRetries {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, mistralAPIBase, bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return "", err
		}
		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return "", readErr
		}

		if resp.StatusCode == http.StatusTooManyRequests && i < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("mistral API error %d: %s", resp.StatusCode, string(respBody))
		}

		var mistralResp mistralResponse
		if err := json.Unmarshal(respBody, &mistralResp); err != nil {
			return "", fmt.Errorf("failed to parse Mistral response: %v", err)
		}
		if len(mistralResp.Choices) == 0 {
			return "", fmt.Errorf("no response from LLM")
		}
		return mistralResp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("max retries exceeded")
}

func (c *MistralClient) GenerateQuiz(ctx context.Context, words []string, progress string, hostLang, targetLang string) ([]Question, error) {
	prompt := fmt.Sprintf(`
You are a language teacher creating an interactive vocabulary quiz.
Generate a quiz for the following target words: %s.

Host Language (Learner's native/UI language): %s
Target Language (Language being learned): %s

Current learner progress:
%s

Instructions:
1. Include two types of questions: "multiple_choice" and "fill_in_the_blank".
2. Ensure a mix of questions. Some should be translation-based (%s to %s or vice versa),
   while others should be entirely in %s (e.g., fill in the blank in a sentence,
   choosing the correct conjugation, or matching a %s definition to a word).
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

The learner's progress is structured into Mastered / Progressing / Struggling concepts. Prioritize words and question types that target Struggling and Progressing concepts. Use Mastered concepts to set baseline difficulty.
Return ONLY a valid JSON array of these question objects. No preamble or markdown formatting.
`, strings.Join(words, ", "), hostLang, targetLang, progress, hostLang, targetLang, targetLang, targetLang)

	text, err := c.complete(ctx, []mistralMessage{{Role: "user", Content: prompt}}, true)
	if err != nil {
		return nil, err
	}

	sanitized := sanitizeJSON(text)
	var questions []Question
	if err := json.Unmarshal([]byte(sanitized), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse quiz JSON: %v", err)
	}
	return questions, nil
}

func (c *MistralClient) UpdateProgress(ctx context.Context, currentProgress string, sessionResults string, randomWords []string) (string, error) {
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

	return c.complete(ctx, []mistralMessage{{Role: "user", Content: prompt}}, false)
}

func (c *MistralClient) Chat(ctx context.Context, progress string, history []Message, userMessage string, hostLang, targetLang string) (*ChatResponse, []Message, error) {
	systemPrompt := fmt.Sprintf(`
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
`, progress, targetLang, hostLang, targetLang)

	messages := []mistralMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range history {
		role := msg.Role
		if role == "model" {
			role = "assistant"
		}
		messages = append(messages, mistralMessage{Role: role, Content: msg.Content})
	}
	messages = append(messages, mistralMessage{Role: "user", Content: userMessage})

	text, err := c.complete(ctx, messages, true)
	if err != nil {
		return nil, nil, err
	}

	sanitized := sanitizeJSON(text)
	var chatRes ChatResponse
	if err := json.Unmarshal([]byte(sanitized), &chatRes); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON response: %v\nResponse was: %s", err, sanitized)
	}

	newHistory := append(history, Message{Role: "user", Content: userMessage})
	newHistory = append(newHistory, Message{Role: "model", Content: chatRes.Response})

	return &chatRes, newHistory, nil
}
