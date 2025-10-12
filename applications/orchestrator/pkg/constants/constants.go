package constants

import "fmt"

const (
	Zero = 0
)

const (
	MessagesTopicPrefix = "messages"
)

func GenerateMessagesTopic(server string) string {
	return fmt.Sprintf("%s.%s", MessagesTopicPrefix, server)
}
