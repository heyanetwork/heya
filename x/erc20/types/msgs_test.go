package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var addr1 = sdk.AccAddress([]byte("addr1_______________")).String()
var addr2 = sdk.AccAddress([]byte("addr2_______________")).String()

func TestMsgTransfer_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   *MsgTransfer
		valid bool
	}{
		{"valid", NewMsgTransfer(addr1, addr2, "1000factory/creator/mytoken"), true},
		{"invalid sender", NewMsgTransfer("bad", addr2, "1000factory/creator/mytoken"), false},
		{"invalid to", NewMsgTransfer(addr1, "bad", "1000factory/creator/mytoken"), false},
		{"invalid amount", NewMsgTransfer(addr1, addr2, "notanumber"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgTransfer_GetSigners(t *testing.T) {
	msg := NewMsgTransfer(addr1, addr2, "100")
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr1, signers[0].String())
}

func TestMsgApprove_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   *MsgApprove
		valid bool
	}{
		{"valid", NewMsgApprove(addr1, addr2, "1000factory/creator/mytoken"), true},
		{"invalid owner", NewMsgApprove("bad", addr2, "1000factory/creator/mytoken"), false},
		{"invalid spender", NewMsgApprove(addr1, "bad", "1000factory/creator/mytoken"), false},
		{"invalid amount", NewMsgApprove(addr1, addr2, "notanumber"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgApprove_GetSigners(t *testing.T) {
	msg := NewMsgApprove(addr1, addr2, "100")
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr1, signers[0].String())
}

func TestMsgTransferFrom_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   *MsgTransferFrom
		valid bool
	}{
		{"valid", NewMsgTransferFrom(addr1, addr2, addr1, "1000factory/creator/mytoken"), true},
		{"invalid caller", NewMsgTransferFrom("bad", addr2, addr1, "1000factory/creator/mytoken"), false},
		{"invalid from", NewMsgTransferFrom(addr1, "bad", addr1, "1000factory/creator/mytoken"), false},
		{"invalid to", NewMsgTransferFrom(addr1, addr2, "bad", "1000factory/creator/mytoken"), false},
		{"invalid amount", NewMsgTransferFrom(addr1, addr2, addr1, "notanumber"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgTransferFrom_GetSigners(t *testing.T) {
	msg := NewMsgTransferFrom(addr1, addr2, addr1, "100")
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr1, signers[0].String())
}

func TestMsgIncreaseAllowance_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   *MsgIncreaseAllowance
		valid bool
	}{
		{"valid", NewMsgIncreaseAllowance(addr1, addr2, "factory/creator/mytoken", "100"), true},
		{"invalid owner", NewMsgIncreaseAllowance("bad", addr2, "factory/creator/mytoken", "100"), false},
		{"invalid spender", NewMsgIncreaseAllowance(addr1, "bad", "factory/creator/mytoken", "100"), false},
		{"zero amount", NewMsgIncreaseAllowance(addr1, addr2, "factory/creator/mytoken", "0"), false},
		{"negative amount", NewMsgIncreaseAllowance(addr1, addr2, "factory/creator/mytoken", "-1"), false},
		{"non-numeric amount", NewMsgIncreaseAllowance(addr1, addr2, "factory/creator/mytoken", "abc"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgIncreaseAllowance_GetSigners(t *testing.T) {
	msg := NewMsgIncreaseAllowance(addr1, addr2, "denom", "100")
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr1, signers[0].String())
}

func TestMsgDecreaseAllowance_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   *MsgDecreaseAllowance
		valid bool
	}{
		{"valid", NewMsgDecreaseAllowance(addr1, addr2, "factory/creator/mytoken", "100"), true},
		{"invalid owner", NewMsgDecreaseAllowance("bad", addr2, "factory/creator/mytoken", "100"), false},
		{"invalid spender", NewMsgDecreaseAllowance(addr1, "bad", "factory/creator/mytoken", "100"), false},
		{"zero amount", NewMsgDecreaseAllowance(addr1, addr2, "factory/creator/mytoken", "0"), false},
		{"non-numeric amount", NewMsgDecreaseAllowance(addr1, addr2, "factory/creator/mytoken", "abc"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgDecreaseAllowance_GetSigners(t *testing.T) {
	msg := NewMsgDecreaseAllowance(addr1, addr2, "denom", "100")
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr1, signers[0].String())
}

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   *MsgUpdateParams
		valid bool
	}{
		{"valid", NewMsgUpdateParams(addr1, addr2), true},
		{"invalid authority", NewMsgUpdateParams("bad", addr2), false},
		{"invalid new authority", NewMsgUpdateParams(addr1, "bad"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	msg := NewMsgUpdateParams(addr1, addr2)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr1, signers[0].String())
}
