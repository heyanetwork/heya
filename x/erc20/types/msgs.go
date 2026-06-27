package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgTransfer(sender, to, amount string) *MsgTransfer {
	return &MsgTransfer{Sender: sender, To: to, Amount: amount}
}

func (m *MsgTransfer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrUnauthorized
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return ErrUnauthorized
	}
	_, err := sdk.ParseCoinNormalized(m.Amount)
	return err
}

func (m *MsgTransfer) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgApprove(owner, spender, amount string) *MsgApprove {
	return &MsgApprove{Owner: owner, Spender: spender, Amount: amount}
}

func (m *MsgApprove) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return ErrUnauthorized
	}
	if _, err := sdk.AccAddressFromBech32(m.Spender); err != nil {
		return ErrUnauthorized
	}
	_, err := sdk.ParseCoinNormalized(m.Amount)
	return err
}

func (m *MsgApprove) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

func NewMsgTransferFrom(caller, from, to, amount string) *MsgTransferFrom {
	return &MsgTransferFrom{Caller: caller, From: from, To: to, Amount: amount}
}

func (m *MsgTransferFrom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Caller); err != nil {
		return ErrUnauthorized
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return ErrUnauthorized
	}
	if _, err := sdk.AccAddressFromBech32(m.To); err != nil {
		return ErrUnauthorized
	}
	_, err := sdk.ParseCoinNormalized(m.Amount)
	return err
}

func (m *MsgTransferFrom) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Caller)
	return []sdk.AccAddress{addr}
}

func NewMsgUpdateParams(authority, newAuthority string) *MsgUpdateParams {
	return &MsgUpdateParams{Authority: authority, NewAuthority: newAuthority}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrUnauthorized.Wrap("invalid authority address")
	}
	if _, err := sdk.AccAddressFromBech32(m.NewAuthority); err != nil {
		return ErrUnauthorized.Wrap("invalid new authority address")
	}
	return nil
}

func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgIncreaseAllowance(owner, spender, denom, amount string) *MsgIncreaseAllowance {
	return &MsgIncreaseAllowance{Owner: owner, Spender: spender, Denom: denom, Amount: amount}
}

func (m *MsgIncreaseAllowance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return ErrUnauthorized
	}
	if _, err := sdk.AccAddressFromBech32(m.Spender); err != nil {
		return ErrUnauthorized
	}
	amt, ok := sdkmath.NewIntFromString(m.Amount)
	if !ok || !amt.IsPositive() {
		return fmt.Errorf("amount must be positive integer")
	}
	return nil
}

func (m *MsgIncreaseAllowance) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

func NewMsgDecreaseAllowance(owner, spender, denom, amount string) *MsgDecreaseAllowance {
	return &MsgDecreaseAllowance{Owner: owner, Spender: spender, Denom: denom, Amount: amount}
}

func (m *MsgDecreaseAllowance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return ErrUnauthorized
	}
	if _, err := sdk.AccAddressFromBech32(m.Spender); err != nil {
		return ErrUnauthorized
	}
	amt, ok := sdkmath.NewIntFromString(m.Amount)
	if !ok || !amt.IsPositive() {
		return fmt.Errorf("amount must be positive integer")
	}
	return nil
}

func (m *MsgDecreaseAllowance) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}
