package event

// Status - event statuses
type Status string

var AllStatuses = map[Status]struct{}{
	StatusDraft:     {},
	StatusPublished: {},
	StatusArchived:  {},
}

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

func (s Status) IsValid() bool {
	_, ok := AllStatuses[s]
	return ok
}

// Type - event types
type Type string

var AllTypes = map[Type]struct{}{
	TypeMeet:    {},
	TypeTracker: {},
}

const (
	TypeMeet    Type = "meet"
	TypeTracker Type = "expense_tracker"
)

func (t Type) IsValid() bool {
	_, ok := AllTypes[t]
	return ok
}
