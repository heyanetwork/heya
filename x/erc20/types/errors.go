package types

import "cosmossdk.io/errors"

var (
	ErrInsufficientAllowance = errors.Register(ModuleName, 1, "insufficient allowance")
	ErrInsufficientBalance   = errors.Register(ModuleName, 2, "insufficient balance")
	ErrUnauthorized          = errors.Register(ModuleName, 3, "unauthorized")
	ErrInvalidParams         = errors.Register(ModuleName, 4, "invalid params")
)
