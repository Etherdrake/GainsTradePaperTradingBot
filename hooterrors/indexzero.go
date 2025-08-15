package hooterrors

import "github.com/pkg/errors"

var ErrIndexZero = errors.New("BTC/USD is the lowest index")
var ErrMaxLeverage = errors.New("100x is the largest leverage")
