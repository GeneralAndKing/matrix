package enum

type Creation uint
type PublishCreationStatus uint
type WhoCanWatch uint

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
const (
	PublishWatch WhoCanWatch = iota + 1
	FriendWatch
	SelfWatch
)
