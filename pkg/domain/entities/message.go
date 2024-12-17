package entities

// CreateMessageRequest the structure for create message request.
type CreateMessageRequest struct {
	TriggerBy string `json:"trigger_by" validate:"required"`
	Qty       int64  `json:"qty" validate:"required,gte=0"`
}

// MessageData the structure for message data.
type MessageData struct {
	Message   string `json:"message"`
	TriggerBy string `json:"trigger_by"`
}

// KafkaTopic for data type string
type KafkaTopic string

const (
	// TopicPublishMessage to consume topic from producer message
	TopicPublishMessage KafkaTopic = "message.publish"
)

// Predefined chatbot responses
var Responses = map[string]string{
	"Hello":          "Hi there! 😊",
	"Weather update": "The weather is sunny and bright! ☀",
	"Tell me a joke": "Why did the chicken cross the road? To get to the other side! 😂",
}

// Fallback response
const FallbackResponse = "I'm sorry, I didn't understand that. 🤔"

// Define Queries
var Queries = []string{
	"Hello",
	"Weather update",
	"Tell me a joke",
	"Good morning",
	"What's your name?",
	"How are you?",
}
