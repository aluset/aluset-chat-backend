package main

import (
  "encoding/json"
  "log"
  "net/http"
  "os"

  "github.com/sashabaranov/go-openai"
)

type ChatRequest struct {
  Message string `json:"message"`
}

type ChatResponse struct {
  Reply string `json:"reply"`
}

func main() {
  client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

  http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
    var chatReq ChatRequest
    err := json.NewDecoder(r.Body).Decode(&chatReq)
    if err != nil {
      http.Error(w, "Invalid request", http.StatusBadRequest)
      return
    }

    resp, err := client.CreateChatCompletion(r.Context(), openai.ChatCompletionRequest{
      Model: "gpt-3.5-turbo",
      Messages: []openai.ChatCompletionMessage{
        {Role: "user", Content: chatReq.Message},
      },
    })
    if err != nil {
      http.Error(w, "Error from OpenAI: "+err.Error(), http.StatusInternalServerError)
      return
    }

    chatResp := ChatResponse{
      Reply: resp.Choices[0].Message.Content,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chatResp)
  })

  log.Println("Server started on :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}