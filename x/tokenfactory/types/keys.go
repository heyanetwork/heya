package types

const (
	ModuleName  = "tokenfactory"
	StoreKey    = ModuleName
	DenomPrefix = "factory"
)

var (
	DenomKeyPrefix       = []byte{0x01}
	SupplyCapKeyPrefix   = []byte{0x02}
	PausedKey             = []byte("paused")
	PendingAdminKeyPrefix = []byte{0x03}
	ParamsKey             = []byte("params")
)

func DenomKey(denom string) []byte {
	return append(DenomKeyPrefix, []byte(denom)...)
}

func SupplyCapKey(denom string) []byte {
	return append(SupplyCapKeyPrefix, []byte(denom)...)
}

func PendingAdminKey(denom string) []byte {
	return append(PendingAdminKeyPrefix, []byte(denom)...)
}
