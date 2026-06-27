package types

import "cosmossdk.io/errors"

var (
	ErrDenomExists     = errors.Register(ModuleName, 1, "denom already exists")
	ErrUnauthorized    = errors.Register(ModuleName, 2, "unauthorized")
	ErrInvalidDenom    = errors.Register(ModuleName, 3, "invalid denom")
	ErrInvalidCreator  = errors.Register(ModuleName, 4, "invalid creator address")
	ErrInvalidSubdenom = errors.Register(ModuleName, 5, "invalid subdenom")
	ErrInsufficientFee = errors.Register(ModuleName, 6, "insufficient denom creation fee")
	ErrDenomNotFound      = errors.Register(ModuleName, 7, "denom not found")
	ErrSupplyCapExceeded   = errors.Register(ModuleName, 8, "supply cap exceeded")
	ErrPaused              = errors.Register(ModuleName, 9, "module is paused")
	ErrPendingAdminNotFound = errors.Register(ModuleName, 10, "pending admin not found")
	ErrNotPendingAdmin     = errors.Register(ModuleName, 11, "only the pending admin can accept")
)
