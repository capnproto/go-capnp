package bbr

import "time"

type sample struct {
	Size     int64
	SendTime time.Time
	AckTime  time.Time
}

func (s sample) DeliveryRate() int64 {
	return s.Size * int64(s.RoundTripTime())
}

func (s sample) RoundTripTime() time.Duration {
	return s.AckTime.Sub(s.SendTime)
}
