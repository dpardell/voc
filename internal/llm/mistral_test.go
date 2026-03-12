package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mistralOKResponse(content string) []byte {
	resp := mistralResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{Message: struct {
				Content string `json:"content"`
			}{Content: content}},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

func TestNewMistralClient(t *testing.T) {
	c := NewMistralClient("key", "")
	if c.modelName != "mistral-small-latest" {
		t.Errorf("expected default model mistral-small-latest, got %s", c.modelName)
	}

	c2 := NewMistralClient("key", "mistral-large-latest")
	if c2.modelName != "mistral-large-latest" {
		t.Errorf("expected model mistral-large-latest, got %s", c2.modelName)
	}
}

func TestMistralGenerateQuiz(t *testing.T) {
	questions := []Question{
		{
			Type:               "multiple_choice",
			Question:           "What does 'chat' mean?",
			Options:            []string{"dog", "cat", "bird", "fish"},
			CorrectAnswerIndex: 1,
			CorrectAnswer:      "cat",
			TargetWord:         "chat",
		},
		{
			Type:          "fill_in_the_blank",
			Question:      "Le ___ est sur le toit.",
			CorrectAnswer: "chat",
			TargetWord:    "chat",
		},
	}
	payload, _ := json.Marshal(questions)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mistralOKResponse(string(payload)))
	}))
	defer srv.Close()

	c := NewMistralClient("test-key", "")
	c.baseURL = srv.URL

	got, err := c.GenerateQuiz(context.Background(), []string{"chat"}, "", "en", "fr")
	if err != nil {
		t.Fatalf("GenerateQuiz returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 questions, got %d", len(got))
	}
	if got[0].Type != "multiple_choice" {
		t.Errorf("expected multiple_choice, got %s", got[0].Type)
	}
	if got[1].CorrectAnswer != "chat" {
		t.Errorf("expected correct_answer 'chat', got %s", got[1].CorrectAnswer)
	}
}

func TestMistralChat(t *testing.T) {
	var capturedBody []byte

	chatPayload := `{"response":"Bonjour!","corrections":[]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(mistralOKResponse(chatPayload))
	}))
	defer srv.Close()

	c := NewMistralClient("test-key", "")
	c.baseURL = srv.URL

	history := []Message{{Role: "user", Content: "Salut"}, {Role: "model", Content: "Salut!"}}
	resp, newHistory, err := c.Chat(context.Background(), "", history, "Comment vas-tu?", "en", "fr")
	if err != nil {
		t.Fatalf("Chat returned error: %v", err)
	}
	if resp.Response != "Bonjour!" {
		t.Errorf("expected response 'Bonjour!', got %s", resp.Response)
	}
	if len(resp.Corrections) != 0 {
		t.Errorf("expected empty corrections, got %v", resp.Corrections)
	}

	// History should have grown by 2 (user + model)
	if len(newHistory) != len(history)+2 {
		t.Errorf("expected history length %d, got %d", len(history)+2, len(newHistory))
	}

	// Verify "model" role is sent as "assistant" in the request
	var req mistralRequest
	if err := json.Unmarshal(capturedBody, &req); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	for _, msg := range req.Messages {
		if msg.Role == "model" {
			t.Error("role 'model' should be sent as 'assistant' to the API")
		}
	}
}

func TestMistralUpdateProgress(t *testing.T) {
	progressMD := "## Summary\nGoing well."

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mistralOKResponse(progressMD))
	}))
	defer srv.Close()

	c := NewMistralClient("test-key", "")
	c.baseURL = srv.URL

	got, err := c.UpdateProgress(context.Background(), "", "session results", []string{"mot"})
	if err != nil {
		t.Fatalf("UpdateProgress returned error: %v", err)
	}
	if got != progressMD {
		t.Errorf("expected %q, got %q", progressMD, got)
	}
}

func TestMistralRetry(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mistralOKResponse("## Progress"))
	}))
	defer srv.Close()

	c := NewMistralClient("test-key", "")
	c.baseURL = srv.URL
	c.retryDelay = time.Millisecond

	got, err := c.UpdateProgress(context.Background(), "", "", nil)
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if got != "## Progress" {
		t.Errorf("unexpected response: %q", got)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestMistralError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
	}))
	defer srv.Close()

	c := NewMistralClient("test-key", "")
	c.baseURL = srv.URL

	_, err := c.UpdateProgress(context.Background(), "", "", nil)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
	errStr := err.Error()
	if len(errStr) == 0 {
		t.Error("expected non-empty error message")
	}
}
