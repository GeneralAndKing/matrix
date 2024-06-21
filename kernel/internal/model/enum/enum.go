package enum

type Work uint
type PublishWorkStatus uint

const (
	PendingWorkStatus PublishWorkStatus = iota
	RunningWorkStatus
	CompletedWorkStatus
	FailedWorkStatus
)

const (
	Video Work = iota + 1
	Teletext
)
