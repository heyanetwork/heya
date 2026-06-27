package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var reValidSubdenom = regexp.MustCompile(`^[a-z0-9_]+$`)

const MaxSubdenomLength = 64

func NewDenom(creator, subdenom string) string {
	return DenomPrefix + "/" + creator + "/" + subdenom
}

func validateSubdenom(subdenom string) error {
	if len(subdenom) == 0 {
		return ErrInvalidSubdenom.Wrap("subdenom cannot be empty")
	}
	if len(subdenom) > MaxSubdenomLength {
		return ErrInvalidSubdenom.Wrapf("subdenom too long, max %d characters", MaxSubdenomLength)
	}
	if !reValidSubdenom.MatchString(subdenom) {
		return ErrInvalidSubdenom.Wrap("subdenom contains invalid characters (allowed: a-z 0-9 _)")
	}
	if subdenom[0] == '.' || subdenom[0] == '-' || subdenom[0] == '_' {
		return ErrInvalidSubdenom.Wrap("subdenom cannot start with a special character")
	}
	if subdenom[len(subdenom)-1] == '.' || subdenom[len(subdenom)-1] == '-' || subdenom[len(subdenom)-1] == '_' {
		return ErrInvalidSubdenom.Wrap("subdenom cannot end with a special character")
	}
	return nil
}

func validateFactoryDenomInMsg(denom string) error {
	if !strings.HasPrefix(denom, DenomPrefix+"/") {
		return ErrInvalidDenom.Wrapf("denom must start with %s/", DenomPrefix)
	}
	return nil
}

func (m *MsgCreateDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrInvalidCreator
	}
	return validateSubdenom(m.Subdenom)
}

func (m *MsgCreateDenom) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgCreateDenom(sender, subdenom string) *MsgCreateDenom {
	return &MsgCreateDenom{Sender: sender, Subdenom: subdenom}
}

func (m *MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrInvalidCreator
	}
	coin, err := sdk.ParseCoinNormalized(m.Amount)
	if err != nil {
		return err
	}
	if !coin.IsPositive() {
		return fmt.Errorf("amount must be positive")
	}
	if err := validateFactoryDenomInMsg(coin.Denom); err != nil {
		return err
	}
	if m.MintTo != "" {
		if _, err := sdk.AccAddressFromBech32(m.MintTo); err != nil {
			return ErrInvalidCreator
		}
	}
	return nil
}

func (m *MsgMint) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgMint(sender, amount, mintTo string) *MsgMint {
	return &MsgMint{Sender: sender, Amount: amount, MintTo: mintTo}
}

func (m *MsgBurn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrInvalidCreator
	}
	coin, err := sdk.ParseCoinNormalized(m.Amount)
	if err != nil {
		return err
	}
	if !coin.IsPositive() {
		return fmt.Errorf("amount must be positive")
	}
	if err := validateFactoryDenomInMsg(coin.Denom); err != nil {
		return err
	}
	return nil
}

func (m *MsgBurn) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgBurn(sender, amount string) *MsgBurn {
	return &MsgBurn{Sender: sender, Amount: amount}
}

func (m *MsgChangeAdmin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrInvalidCreator
	}
	if _, err := sdk.AccAddressFromBech32(m.NewAdmin); err != nil {
		return ErrInvalidCreator
	}
	if err := validateFactoryDenomInMsg(m.Denom); err != nil {
		return err
	}
	return nil
}

func (m *MsgChangeAdmin) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgChangeAdmin(sender, denom, newAdmin string) *MsgChangeAdmin {
	return &MsgChangeAdmin{Sender: sender, Denom: denom, NewAdmin: newAdmin}
}

func (m *MsgForceTransfer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrInvalidCreator
	}
	if _, err := sdk.AccAddressFromBech32(m.DestAddr); err != nil {
		return ErrInvalidCreator
	}
	coin, err := sdk.ParseCoinNormalized(m.Amount)
	if err != nil {
		return err
	}
	if !coin.IsPositive() {
		return fmt.Errorf("amount must be positive")
	}
	if err := validateFactoryDenomInMsg(coin.Denom); err != nil {
		return err
	}
	return nil
}

func (m *MsgForceTransfer) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgForceTransfer(sender, amount, destAddr string) *MsgForceTransfer {
	return &MsgForceTransfer{Sender: sender, Amount: amount, DestAddr: destAddr}
}

func NewMsgForceTransferFull(sender, amount, destAddr, fromAddress string) *MsgForceTransfer {
	return &MsgForceTransfer{Sender: sender, Amount: amount, DestAddr: destAddr, FromAddress: fromAddress}
}

func (m *MsgAcceptAdmin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Admin); err != nil {
		return ErrInvalidCreator
	}
	if err := validateFactoryDenomInMsg(m.Denom); err != nil {
		return err
	}
	return nil
}

func (m *MsgAcceptAdmin) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Admin)
	return []sdk.AccAddress{addr}
}

func NewMsgAcceptAdmin(admin, denom string) *MsgAcceptAdmin {
	return &MsgAcceptAdmin{Admin: admin, Denom: denom}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidCreator.Wrap("invalid authority address")
	}
	coin, err := sdk.ParseCoinNormalized(m.DenomCreationFee)
	if err != nil {
		return ErrInsufficientFee.Wrap("invalid denom creation fee")
	}
	params := Params{DenomCreationFee: coin}
	return params.Validate()
}

func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgUpdateParams(authority, denomCreationFee string) *MsgUpdateParams {
	return &MsgUpdateParams{Authority: authority, DenomCreationFee: denomCreationFee}
}
