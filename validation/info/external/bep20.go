package external

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/trustwallet/assets-go-libs/http"
)

var (
	holdersRegexp  = regexp.MustCompile(`(\d+)\saddresses`)
	decimalsRegexp = regexp.MustCompile(`(\d+)\s<\/div>`)
	symbolRegexp   = regexp.MustCompile(`<b>(\w+)<\/b>\s<span`)
)

func GetTokenInfoForBEP20(tokenID string) (*TokenInfo, error) {
	url := fmt.Sprintf("https://bscscan.com/token/%s", tokenID)

	data, err := http.GetHTTPResponseBytes(url)
	if err != nil {
		return nil, err
	}

	// Remove all "," from content.
	pageContent := strings.ReplaceAll(string(data), ",", "")

	var holders, decimals int
	var symbol string

	match := symbolRegexp.FindStringSubmatch(pageContent)
	if len(match) > 1 {
		symbol = match[1]
		if err != nil {
			return nil, err
		}
	}

	match = decimalsRegexp.FindStringSubmatch(pageContent)
	if len(match) > 1 {
		decimals, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
	}

	match = holdersRegexp.FindStringSubmatch(pageContent)
	if len(match) > 1 {
		holders, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
	}

	return &TokenInfo{
		Symbol:       symbol,
		Decimals:     decimals,
		HoldersCount: holders,
	}, nil
}
