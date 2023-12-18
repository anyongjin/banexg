package binance

import (
	"github.com/anyongjin/banexg"
	"github.com/anyongjin/banexg/utils"
	"strconv"
)

func (mar *BnbMarket) GetPrecision() *banexg.Precision {
	var pre = banexg.Precision{}
	if mar.QuantityPrecision > 0 {
		pre.Amount = mar.QuantityPrecision
	} else if mar.QuantityScale > 0 {
		pre.Amount = mar.QuantityScale
	}
	if mar.PricePrecision > 0 {
		pre.Price = mar.PricePrecision
	} else if mar.PriceScale > 0 {
		pre.Price = mar.PriceScale
	}
	pre.Base = mar.BaseAssetPrecision
	pre.Quote = mar.QuotePrecision
	return &pre
}

func (mar *BnbMarket) GetMarketLimits() (*banexg.MarketLimits, int, int) {
	minQty, _ := strconv.ParseFloat(mar.MinQty, 64)
	maxQty, _ := strconv.ParseFloat(mar.MaxQty, 64)
	var filters = make(map[string]BnbFilter)
	for _, flt := range mar.Filters {
		filters[utils.GetMapVal(flt, "filterType", "")] = flt
	}
	var res = banexg.MarketLimits{
		Amount: &banexg.LimitRange{
			Min: minQty,
			Max: maxQty,
		},
		Leverage: &banexg.LimitRange{},
		Price:    &banexg.LimitRange{},
		Cost:     &banexg.LimitRange{},
		Market:   &banexg.LimitRange{},
	}
	var pricePrec, amountPrec int
	if flt, ok := filters["PRICE_FILTER"]; ok {
		// PRICE_FILTER reports zero values for maxPrice
		// since they updated filter types in November 2018
		// https://github.com/ccxt/ccxt/issues/4286
		// therefore limits['price']['max'] doesn't have any meaningful value except None
		res.Price.Min = utils.GetMapFloat(flt, "minPrice")
		res.Price.Max = utils.GetMapFloat(flt, "maxPrice")
		precText := utils.GetMapVal(flt, "tickSize", "")
		pricePrec = utils.PrecisionFromString(precText)
	}
	if flt, ok := filters["LOT_SIZE"]; ok {
		res.Amount.Min = utils.GetMapFloat(flt, "minQty")
		res.Amount.Max = utils.GetMapFloat(flt, "maxQty")
		amtText := utils.GetMapVal(flt, "stepSize", "")
		amountPrec = utils.PrecisionFromString(amtText)
	}
	if flt, ok := filters["MARKET_LOT_SIZE"]; ok {
		res.Market.Min = utils.GetMapFloat(flt, "minQty")
		res.Market.Max = utils.GetMapFloat(flt, "maxQty")
	}
	if flt, ok := filters["MIN_NOTIONAL"]; ok {
		res.Cost.Min = utils.GetMapFloat(flt, "notional")
	} else if flt, ok := filters["NOTIONAL"]; ok {
		res.Cost.Min = utils.GetMapFloat(flt, "minNotional")
		res.Cost.Max = utils.GetMapFloat(flt, "maxNotional")
	}
	return &res, pricePrec, amountPrec
}

func (a BnbAsset) ToStdAsset(e *banexg.Exchange) *banexg.Asset {
	free, _ := strconv.ParseFloat(a.Free, 64)
	lock, _ := strconv.ParseFloat(a.Locked, 64)
	borr, _ := strconv.ParseFloat(a.Borrowed, 64)
	inst, _ := strconv.ParseFloat(a.Interest, 64)
	code := e.SafeCurrencyCode(a.Asset)
	return &banexg.Asset{
		Code:  code,
		Free:  free,
		Used:  lock,
		Total: lock + free,
		Debt:  borr + inst,
	}
}

func (a *FutureAsset) ToStdAsset(e *banexg.Exchange) *banexg.Asset {
	code := e.SafeCurrencyCode(a.Asset)
	free, _ := strconv.ParseFloat(a.AvailableBalance, 64)
	used, _ := strconv.ParseFloat(a.InitialMargin, 64)
	total, _ := strconv.ParseFloat(a.MarginBalance, 64)
	return &banexg.Asset{
		Code:  code,
		Free:  free,
		Used:  used,
		Total: total,
	}
}
