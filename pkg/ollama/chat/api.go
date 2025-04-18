package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	RoleUser string = "user"
	RoleTool string = "tool"
)

type Chat interface {
	Send(ctx context.Context, input string, role string) (Response, error)
}

type ChatPrompt struct {
	Model    string                 `json:"model"`
	Messages []Message              `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options"`
	Tools    []Tool                 `json:"tools"`
}

type Response struct {
	Msg Message `json:"message"`
}

type Message struct {
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	ToolCalls []map[string]ToolCalls `json:"tool_calls"`
}

type ToolCalls struct {
	Name string
	//TODO: Fix me
	Arguments json.RawMessage
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type     string              `json:"type"`
	Property map[string]Property `json:"property"`
	Required []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum"`
}

// client implements the Chat interface for interacting with the Ollama API.
// It handles constructing requests, sending messages, and parsing responses.
type client struct {
	// Url is the Ollama API endpoint
	Url string
	// Model is the name of the language model to use
	Model string
}

// OpenClient creates a new client instance with the specified API URL and model name.
// It returns the client as a Chat interface implementation.
func OpenClient(url, model string) (Chat, error) {
	return &client{
		Url:   url,
		Model: model,
	}, nil
}

// Send submits a chat message to the Ollama API with the specified input text and role.
// It uses a synchronous (non-streaming) request and includes a predefined weather tool (this will change).
// The context can be used to control cancellation.
//
// Parameters:
//   - ctx: Context for the request
//   - input: The message content to send
//   - role: The role for the message (e.g., "user", "assistant", "tool")
func (c *client) Send(ctx context.Context, input string, role string) (Response, error) {
	output := Response{}

	gr := ChatPrompt{
		Model:  c.Model,
		Stream: false,
		Messages: []Message{
			{
				Role:    role,
				Content: input,
			},
		},
		Options: nil,
		Tools: []Tool{{
			//TODO pull in from the pointer receiveer
			Type: "function",
			Function: Function{
				Name:        "get_current_weather",
				Description: "Get the current weather for a location",
				Parameters: Parameters{
					Type: "object",
					Property: map[string]Property{
						"location": {
							Type:        "string",
							Description: "The location to get the weather for, e.g. San Francisco, CA",
						},
					},
					Required: []string{
						"location",
					},
				},
			},
		}},
	}

	//TODO: Fix error handling
	b, _ := json.Marshal(gr)
	bodyReader := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, c.Url, bodyReader)
	if err != nil {
		//TODO: Fix me
		return output, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: Fix error handling
		return output, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		raw, err := io.ReadAll(res.Body)
		if err != nil {
			//TODO: Fix error handling
			return output, err
		}
		uErr := json.Unmarshal(raw, &output)
		if uErr != nil {
			//TODO: Fix error handling
			return output, uErr
		}
	}
	return output, nil
}
