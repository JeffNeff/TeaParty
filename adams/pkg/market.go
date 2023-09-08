// fetch marketplace data from the API
package be

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	gecko "github.com/superoo7/go-gecko/v3"
	geckoTypes "github.com/superoo7/go-gecko/v3/types"
)

func FetchMarketPriceInUSD(coin string) (*big.Int, error) {
	if coin == GRAMS {
		bignum, ok := new(big.Int).SetString("1000000000000000000", 0)
		if !ok {
			return nil, fmt.Errorf("creating new mo static thinggy")
		}
		return bignum, nil
	}

	if coin == BSCUSDT {
		bignum, ok := new(big.Int).SetString("1000000000000000000", 0)
		if !ok {
			return nil, fmt.Errorf("creating new mo static thinggy")
		}
		return bignum, nil
	}

	if coin == POL {
		coin = "matic-network"
	}

	cg := gecko.NewClient(nil)
	// find specific coins
	vsCurrency := "usd"
	ids := []string{coin}
	perPage := 1
	page := 1
	sparkline := true
	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{pcp.PCP1h, pcp.PCP24h, pcp.PCP7d, pcp.PCP14d, pcp.PCP30d, pcp.PCP200d, pcp.PCP1y}
	order := geckoTypes.OrderTypeObject.MarketCapDesc
	market, err := cg.CoinsMarket(vsCurrency, ids, order, perPage, page, sparkline, priceChangePercentage)
	if err != nil {
		return big.NewInt(0), err
	}

	fmt.Println("price in USD: ", (*market)[0].CurrentPrice)
	bi := bigIntViaString((*market)[0].CurrentPrice)
	bi.Div(bi, big.NewInt(100))
	fmt.Println("price in USD: ", bi.String())
	return bi, nil
}

func bigIntViaString(flt float64) (b *big.Int) {

	if math.IsNaN(flt) || math.IsInf(flt, 0) {
		return nil // illegal case
	}

	var in = strconv.FormatFloat(flt, 'f', -1, 64)

	const parts = 2

	var ss = strings.SplitN(in, ".", parts)

	// protect from numbers without period
	if len(ss) != parts {
		ss = append(ss, "0")
	}

	// protect from ".0" and "0." values
	if ss[0] == "" {
		ss[0] = "0"
	}

	if ss[1] == "" {
		ss[1] = "0"
	}

	const (
		base     = 10
		fraction = 20
	)

	// get fraction length
	var fract = len(ss[1])
	if fract > fraction {
		ss[1], fract = ss[1][:fraction], fraction
	}

	in = strings.Join([]string{ss[0], ss[1]}, "")
	// convert to big integer from the string
	b, _ = big.NewInt(0).SetString(in, base)
	if fract == fraction {
		return // ready
	}
	// fract < 20, * (20 - fract)
	var (
		ten = big.NewInt(base)
		exp = ten.Exp(ten, big.NewInt(fraction-int64(fract)), nil)
	)
	b = b.Mul(b, exp)
	return

}
