package kucoinapi

// kucoin SPOT Symbols REST 获取交易对列表 (NONE)
func (client *SpotRestClient) NewSymbols() *SpotSymbolsApi {
	return &SpotSymbolsApi{
		client: client,
		req:    &SpotSymbolsReq{},
	}
}
func (api *SpotSymbolsApi) Do() (*SpotSymbolsRes, error) {
	url := kucoinHandlerRequestApiWithPathQueryParam(SPOT, api.req, SpotApiMap[SpotSymbols])
	return kucoinCallApi[SpotSymbolsRes](api.client.c, url, NIL_REQBODY, GET)
}

// kucoin SPOT timestamp REST 服务器时间 (NONE)
func (client *SpotRestClient) NewTimestamp() *SpotTimestampApi {
	return &SpotTimestampApi{
		client: client,
		req:    &SpotTimestampReq{},
	}
}
func (api *SpotTimestampApi) Do() (*SpotTimestampRes, error) {
	url := kucoinHandlerRequestApiWithPathQueryParam(SPOT, api.req, SpotApiMap[SpotTimestamp])
	return kucoinCallApi[SpotTimestampRes](api.client.c, url, NIL_REQBODY, kucoinHandlerReq(api.req))
}
