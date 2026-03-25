package domain

import "errors"

var ErrValidationFailed = errors.New("validation failed")
var ErrFraudDetected = errors.New("fraud detected")
var ErrRegionNotFound = errors.New("no matching region")