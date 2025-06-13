package base

import "fmt"

type Message struct {
	Role    RoleType
	Content string
}

func (m *Message) toString() string {
	return fmt.Sprintf("角色: %s\n内容: %s\n", m.Role, m.Content)
}

type RoleType string

const (
	// Assistant is the role of an assistant, means the message is returned by ChatModel.
	Assistant RoleType = "assistant"
	// User is the role of a user, means the message is a user message.
	User RoleType = "user"
	// System is the role of a system, means the message is a system message.
	System RoleType = "system"
	// Tool is the role of a tool, means the message is a tool call output.
	Tool RoleType = "tool"
)

func SystemMessage(content string) *Message {
	return &Message{
		Role:    System,
		Content: content,
	}
}

func UserMessage(content string) *Message {
	return &Message{
		Role:    User,
		Content: content,
	}
}

func AssistantMessage(content string) *Message {
	return &Message{
		Role:    Assistant,
		Content: content,
	}
}

func ToolMessage(content string) *Message {
	return &Message{
		Role:    Tool,
		Content: content,
	}
}
