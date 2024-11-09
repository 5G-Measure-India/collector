package util

type Type = chan struct{}

var (
	Stop          = make(Type)
	GeoDone       = make(Type)
	PhyDone       = make(Type)
	PingDone      = make(Type)
	SpeedtestDone = make(Type)
)
