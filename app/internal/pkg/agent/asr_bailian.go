package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/solikewind/happyeat/app/internal/config"
)

const (
	defaultASRModel          = "paraformer-realtime-v2"
	defaultTranscribePath    = "/asr/transcribe"
	defaultHotwordCreatePath = "/asr/hotwords"
)

type BailianASRClient struct {
	baseURL           string
	apiKey            string
	model             string
	transcribePath    string
	hotwordCreatePath string
	httpClient        *http.Client
}

type TranscribeInput struct {
	AudioBase64 string   `json:"audio_base64"`
	Format      string   `json:"format,omitempty"`      // wav/mp3/pcm...
	SampleRate  int      `json:"sample_rate,omitempty"` // 16000
	HotwordIDs  []string `json:"hotword_ids,omitempty"`
}

type TranscribeResult struct {
	Text string
	Raw  json.RawMessage
}

type CreateHotwordInput struct {
	Name        string   `json:"name"`
	Words       []string `json:"words"`
	Description string   `json:"description,omitempty"`
}

type CreateHotwordResult struct {
	ID   string
	Raw  json.RawMessage
	Meta map[string]any
}

func NewBailianASRClient(c config.ASR) (*BailianASRClient, error) {
	apiKey := strings.TrimSpace(c.APIKey)
	baseURL := strings.TrimSpace(c.BaseURL)
	if apiKey == "" || baseURL == "" {
		return nil, errors.New("asr api key or base url is empty")
	}

	model := strings.TrimSpace(c.Model)
	if model == "" {
		model = defaultASRModel
	}

	transcribePath := strings.TrimSpace(c.TranscribePath)
	if transcribePath == "" {
		transcribePath = defaultTranscribePath
	}

	hotwordCreatePath := strings.TrimSpace(c.HotwordCreatePath)
	if hotwordCreatePath == "" {
		hotwordCreatePath = defaultHotwordCreatePath
	}

	return &BailianASRClient{
		baseURL:           strings.TrimRight(baseURL, "/"),
		apiKey:            apiKey,
		model:             model,
		transcribePath:    transcribePath,
		hotwordCreatePath: hotwordCreatePath,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}, nil
}

func (c *BailianASRClient) Transcribe(ctx context.Context, in TranscribeInput) (*TranscribeResult, error) {
	if strings.TrimSpace(in.AudioBase64) == "" {
		return nil, errors.New("audio_base64 is empty")
	}

	payload := map[string]any{
		"model": c.model,
		"input": map[string]any{
			"audio_base64": in.AudioBase64,
			"format":       in.Format,
			"sample_rate":  in.SampleRate,
			"hotword_ids":  in.HotwordIDs,
		},
	}

	raw, err := c.postJSON(ctx, c.transcribePath, payload)
	if err != nil {
		return nil, err
	}

	text := extractTextFromAnyJSON(raw)
	return &TranscribeResult{
		Text: text,
		Raw:  raw,
	}, nil
}

func (c *BailianASRClient) CreateHotword(ctx context.Context, in CreateHotwordInput) (*CreateHotwordResult, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, errors.New("hotword name is empty")
	}
	if len(in.Words) == 0 {
		return nil, errors.New("hotword words are empty")
	}

	payload := map[string]any{
		"name":        in.Name,
		"words":       in.Words,
		"description": in.Description,
	}

	raw, err := c.postJSON(ctx, c.hotwordCreatePath, payload)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	_ = json.Unmarshal(raw, &m)
	id := findFirstString(m, "id", "hotword_id", "data.id", "output.id")

	return &CreateHotwordResult{
		ID:   id,
		Raw:  raw,
		Meta: m,
	}, nil
}

func (c *BailianASRClient) postJSON(ctx context.Context, path string, payload any) (json.RawMessage, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	url := c.baseURL + ensureLeadingSlash(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call asr api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read asr response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("asr api failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return body, nil
}

func ensureLeadingSlash(s string) string {
	if strings.HasPrefix(s, "/") {
		return s
	}
	return "/" + s
}

func extractTextFromAnyJSON(raw json.RawMessage) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	return findFirstString(m, "text", "output.text", "result.text", "data.text")
}

func findFirstString(root map[string]any, paths ...string) string {
	for _, p := range paths {
		if v, ok := digValue(root, p); ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func digValue(root map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var cur any = root
	for _, part := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		v, ok := m[part]
		if !ok {
			return nil, false
		}
		cur = v
	}
	return cur, true
}
