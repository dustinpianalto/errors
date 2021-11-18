package errors

type Username string
type Method string
type Message string
type Kind uint16

const (
	Other      Kind = iota // Unknown error or something that doesn't fit other categories
	Internal               // Internal error that should not be shown to user
	Invalid                // Operation is not permitted for this type of item
	Incorrect              // Incorrect configuration or values
	Permission             // Permission denied
	IO                     // External IO error
	Conflict               // The item already exists
	NotFound               // The item does not exist
	Malformed              // The request format is not valid
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Internal:
		return "internal error"
	case Invalid:
		return "invalid operation"
	case Incorrect:
		return "incorrect configuration"
	case Permission:
		return "permission denied"
	case IO:
		return "I/O error"
	case Conflict:
		return "item already exists"
	case NotFound:
		return "item does not exist"
	case Malformed:
		return "malformed request"
	}
	return "unknown type"
}
