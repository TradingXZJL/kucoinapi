package kucoinapi

type SpotApi int

const (
	// 通用接口
	SpotSymbols   = iota // GET 获取交易对列表
	SpotTimestamp        // GET 服务器时间
)

var SpotApiMap = map[SpotApi]string{
	// 通用接口
	SpotSymbols:   "/api/v2/symbols",   // GET 获取交易对列表
	SpotTimestamp: "/api/v1/timestamp", // GET 服务器时间
}
