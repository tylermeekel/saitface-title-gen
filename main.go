package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sashabaranov/go-openai"
)

//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

type Server struct {
	OAIClient *openai.Client
}

func main() {
	server := &Server{}

	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	server.OAIClient = client

	mux := chi.NewMux()

	mux.Post("/", server.HandleGenerateTitle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, mux)
}

func (s *Server) GetTitle(interests []string) string {

	str := strings.Join(interests, "and ")
	str = "Write 1 interesting question related to the topics " + str + " that can attract people to engage in fun, meaningful discussions based on them."

	resp, err := s.OAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: str,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func (s *Server) HandleGenerateTitle(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Interests []string `json:"interests"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	content := s.GetTitle(req.Interests)

	response := struct {
		Title string `json:"title"`
	}{
		Title: content,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
