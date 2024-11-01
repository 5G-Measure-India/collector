package channel

type Type = chan struct{}

var (
	Stop          = make(Type)
	PhyDone       = make(Type)
	PingDone      = make(Type)
	SpeedtestDone = make(Type)
)

func StopAll() {
	close(Stop)

	<-SpeedtestDone
	<-PingDone
	<-PhyDone
}
