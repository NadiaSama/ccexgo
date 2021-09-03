package exchange

import "time"

type (
	KlineResolution int

	Kline struct {
		Symbol Symbol
		Open   float64
		Close  float64
		High   float64
		Low    float64
		Volume float64
		Time   time.Time
		Raw    interface{}
	}

	KlineReq struct {
		Symbol     Symbol
		StartTime  time.Time
		EndTime    time.Time
		Limit      int
		Resolution KlineResolution
	}
)

const (
	KlineResolution1m KlineResolution = iota
	KlineResolution5m
	KlineResolution15m
	KlineResolution30m
	KlineResolution1h
	KlineResolution4h
	KlineResolution1D
	KlineResolution1W
)

var (
	resolutionMap = map[KlineResolution]string{
		KlineResolution1m:  "1m",
		KlineResolution5m:  "5m",
		KlineResolution15m: "15m",
		KlineResolution30m: "30m",
		KlineResolution1h:  "1h",
		KlineResolution4h:  "4h",
		KlineResolution1D:  "1d",
		KlineResolution1W:  "1w",
	}

	resolutionTS = map[KlineResolution]int{
		KlineResolution1m:  60,
		KlineResolution5m:  300,
		KlineResolution15m: 900,
		KlineResolution30m: 1800,
		KlineResolution1h:  3600,
		KlineResolution4h:  14400,
		KlineResolution1D:  86400,
		KlineResolution1W:  604800,
	}
)

func (kr KlineResolution) Secs() int {
	return resolutionTS[kr]
}

func (kr KlineResolution) String() string {
	return resolutionMap[kr]
}

func NewKlineReq(symbol Symbol, resolution KlineResolution) *KlineReq {
	return &KlineReq{
		Symbol:     symbol,
		Resolution: resolution,
	}
}

func (kr *KlineReq) SetLimit(l int) *KlineReq {
	kr.Limit = l
	return kr
}

func (kr *KlineReq) SetStartTime(st time.Time) *KlineReq {
	kr.StartTime = st
	return kr
}

func (kr *KlineReq) SetEndTime(et time.Time) *KlineReq {
	kr.EndTime = et
	return kr
}
