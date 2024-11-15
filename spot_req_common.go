package kucoinapi

type SpotSymbolsReq struct{}
type SpotSymbolsApi struct {
	client *SpotRestClient
	req    *SpotSymbolsReq
}

type SpotTimestampReq struct{}
type SpotTimestampApi struct {
	client *SpotRestClient
	req    *SpotTimestampReq
}
