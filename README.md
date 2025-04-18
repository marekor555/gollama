# 🦙 gollama – Go Client for the Ollama API

`gollama` is a lightweight and efficient Go client for interacting with [Ollama](https://ollama.com)'s local language models.  
It supports both single-prompt generation and multi-turn chat interactions through a simple API.

---

## 📦 Installation

```bash
go get github.com/marekor555/gollama
```

---

## 🔧 Quick Start

### Create a Model

```go
model, err := gollama.CreateModel("llama2", "http://localhost:11434", "You are a helpful assistant.")
if err != nil {
    log.Fatal(err)
}
```

### Generate Text from a Prompt

```go
response, err := model.Generate("Write a short poem about Go.")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response)
```

---

## 💬 Chat Interface

The chat interface keeps track of the conversation history and enables contextual interactions.

### Create a Chat Instance

```go
chat, err := gollama.CreateChat("llama2", "http://localhost:11434", "You are a helpful assistant.")
if err != nil {
    log.Fatal(err)
}
```

### Send and Receive Messages

```go
err = chat.Send("What's the capital of France?")
if err != nil {
    log.Fatal(err)
}

response, err := chat.Receive()
if err != nil {
    log.Fatal(err)
}
fmt.Println(response)
```

### Shortcut for One-Turn Interaction

```go
response, err := chat.SendAndReceive("Explain recursion.")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response)
```

---

## 📚 API Overview

### `Model` Struct

Represents a connection to an Ollama model.

#### Fields
- `Name string` – Model name (e.g., `"llama2"`)
- `Addr string` – Address of the Ollama server
- `System string` – System's prompt to condition the model

#### Methods

- `ListModels() ([]string, error)`  
  Returns all available model names from the server.

- `Generate(msg string) (string, error)`  
  Sends a prompt to the model and returns a generated response.

- `changeModel(name string) error`  
  Changes the current model if it exists on the server.

---

### `Chat` Struct

Represents a conversational context with the model.

#### Fields
- `Model Model` – The underlying model used for responses
- `Messages []Message` – History of chat messages

#### Methods

- `Send(msg string) error`  
  Adds a user message to the chat context.

- `Receive() (string, error)`  
  Sends the full conversation and receives a model-generated reply.

- `SendAndReceive(msg string) (string, error)`  
  Sends a message and returns the assistant's reply in one call.

---

## 🧪 Example Output

```go
model.Generate("What's the meaning of life?")
// resp: "The meaning of life is a philosophical question with many interpretations..."
```

```go
chat.SendAndReceive("Tell me a programming joke.")
// resp: "Why do programmers prefer dark mode? Because light attracts bugs."
```