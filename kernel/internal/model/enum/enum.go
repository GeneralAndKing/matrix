package enum

type Creation uint
type PublishCreationStatus uint

const (
	PendingCreationStatus PublishCreationStatus = iota + 1
	RunningCreationStatus
	CompletedCreationStatus
	FailedCreationStatus
)

const (
	Video Creation = iota + 1
	Teletext
)
