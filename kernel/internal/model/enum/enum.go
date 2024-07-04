package enum

type Creation uint
type PublishCreationStatus uint
type WhoCanWatch uint
type LiveMonitorStatus uint

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
const (
	PauseLiveMonitorStatus LiveMonitorStatus = iota + 1
	RunningLiveMonitorStatus
	NotExistLiveMonitorStatus
)
