package kucoinapi

type SpotTimestampRes int

type SpotSymbolsResRow struct {
	Symbol          string `json:"symbol"`          // 交易對唯一標識碼，重命名後不會改變
	Name            string `json:"name"`            // 交易對名稱，重命名後會改變
	BaseCurrency    string `json:"baseCurrency"`    // 商品貨幣，指一個交易對的交易對象，即寫在靠前部分的資產名
	QuoteCurrency   string `json:"quoteCurrency"`   // 計價幣種，指一個交易對的定價資產，即寫在靠後部分資產名
	Market          string `json:"market"`          // 交易市場
	BaseMinSize     string `json:"baseMinSize"`     // 下單時size的最小值
	QuoteMinSize    string `json:"quoteMinSize"`    // 下市價單，funds的最小值
	BaseMaxSize     string `json:"baseMaxSize"`     // 下單，size的最大值
	QuoteMaxSize    string `json:"quoteMaxSize"`    // 下市價單，funds的最大值
	BaseIncrement   string `json:"baseIncrement"`   // 數量增量，下單的數量增量必須為正整數倍，這裡的 size 指的是下單的基礎貨幣數量。例如，對於 ETH-USDT交易對，若baseIncrement=0.0000001，則下單數量可以是 1.0000001，但不可以是 1.00000001
	QuoteIncrement  string `json:"quoteIncrement"`  // 市價單的資金增量，下單的資金數量必須為資金增量的正整數倍。例如，對於 ETH-USDT 交易對，資金（funds）為報價貨幣（quoteCurrency），若資金增量（quoteIncrement）為 0.000001，則資金的 USDT 數量可以是 3000.000001，但不可以是 3000.0000001
	PriceIncrement  string `json:"priceIncrement"`  // 價格增量，下單的價格必須為價格增量的正整數倍。例如，對於 ETH-USDT 交易對，若價格增量（priceIncrement）為 0.01，則下單價格可以是 3000.01，但不可以是 3000.001
	FeeCurrency     string `json:"feeCurrency"`     // 交易計算手續費的幣種
	EnableTrading   bool   `json:"enableTrading"`   // 是否可以用於交易
	IsMarginEnabled bool   `json:"isMarginEnabled"` // 是否支持槓桿
	PriceLimitRate  string `json:"priceLimitRate"`  // 價格保護閾值
	MinFunds        string `json:"minFunds"`        // 最小交易金額
}
type SpotSymbolsRes []SpotSymbolsResRow
