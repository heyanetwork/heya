package types

const (
	ModuleName = "erc20"
	StoreKey   = ModuleName
)

var (
	AllowanceKeyPrefix = []byte{0x01}
	ParamsKey          = []byte("params")
)

func AllowanceKey(owner, spender, denom string) []byte {
	return append(append(AllowanceKeyPrefix, []byte(owner+"/"+spender+"/"+denom)...), 0)
}
