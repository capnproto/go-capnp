package bbr

type stateName int

const (
	probeBWState stateName = iota
	probeRTTState
	startupState
	drainState
)
