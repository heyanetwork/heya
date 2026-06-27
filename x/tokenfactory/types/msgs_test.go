package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"heya/x/tokenfactory/types"
)

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("heya", "heyapub")
}

func validSenderAddr(t *testing.T) string {
	t.Helper()
	addr := sdk.AccAddress(make([]byte, 20))
	bech, err := bech32.ConvertAndEncode("heya", addr)
	if err != nil {
		t.Fatal(err)
	}
	return bech
}

func otherSenderAddr(t *testing.T) string {
	t.Helper()
	addr := sdk.AccAddress(make([]byte, 20))
	addr[0] = 1
	bech, err := bech32.ConvertAndEncode("heya", addr)
	if err != nil {
		t.Fatal(err)
	}
	return bech
}

func TestValidateSubdenom(t *testing.T) {
	sender := validSenderAddr(t)

	tests := []struct {
		name      string
		subdenom  string
		expectErr bool
	}{
		{"valid simple", "mytoken", false},
		{"valid with numbers", "token123", false},
		{"hyphens not allowed", "my-token", true},
		{"valid with underscores", "my_token", false},
		{"dots not allowed", "my.token", true},
		{"uppercase not allowed", "MyToken", true},
		{"empty", "", true},
		{"too long", string(make([]byte, 65)), true},
		{"special chars", "token@#$\x00", true},
		{"space", "my token", true},
		{"starts with dot", ".token", true},
		{"starts with hyphen", "-token", true},
		{"starts with underscore", "_token", true},
		{"ends with dot", "token.", true},
		{"ends with hyphen", "token-", true},
		{"ends with underscore", "token_", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgCreateDenom(sender, tt.subdenom)
			err := msg.ValidateBasic()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMsgCreateDenom_ValidateBasic(t *testing.T) {
	sender := validSenderAddr(t)

	tests := []struct {
		name      string
		sender    string
		subdenom  string
		expectErr bool
	}{
		{"valid", sender, "mytoken", false},
		{"invalid sender", "invalid", "mytoken", true},
		{"empty subdenom", sender, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgCreateDenom(tt.sender, tt.subdenom)
			err := msg.ValidateBasic()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMsgMint_ValidateBasic(t *testing.T) {
	sender := validSenderAddr(t)
	denom := "factory/" + sender + "/mytoken"

	tests := []struct {
		name      string
		sender    string
		amount    string
		mintTo    string
		expectErr bool
	}{
		{"valid", sender, "1000" + denom, "", false},
		{"valid with mint_to", sender, "1000" + denom, sender, false},
		{"invalid sender", "bad", "1000" + denom, "", true},
		{"invalid amount", sender, "bad", "", true},
		{"zero amount", sender, "0" + denom, "", true},
		{"negative amount", sender, "-1" + denom, "", true},
		{"non-factory denom", sender, "1000uheya", "", true},
		{"missing factory prefix", sender, "1000" + sender + "/mytoken", "", true},
		{"invalid mint_to", sender, "1000" + denom, "bad", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgMint(tt.sender, tt.amount, tt.mintTo)
			err := msg.ValidateBasic()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMsgBurn_ValidateBasic(t *testing.T) {
	sender := validSenderAddr(t)
	denom := "factory/" + sender + "/mytoken"

	tests := []struct {
		name      string
		sender    string
		amount    string
		expectErr bool
	}{
		{"valid", sender, "1000" + denom, false},
		{"invalid sender", "bad", "1000" + denom, true},
		{"invalid amount", sender, "bad", true},
		{"zero amount", sender, "0" + denom, true},
		{"non-factory denom", sender, "1000uheya", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgBurn(tt.sender, tt.amount)
			err := msg.ValidateBasic()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMsgChangeAdmin_ValidateBasic(t *testing.T) {
	sender := validSenderAddr(t)
	other := otherSenderAddr(t)

	tests := []struct {
		name      string
		sender    string
		denom     string
		newAdmin  string
		expectErr bool
	}{
		{"valid", sender, "factory/" + sender + "/mytoken", other, false},
		{"invalid sender", "bad", "factory/" + sender + "/mytoken", other, true},
		{"invalid new_admin", sender, "factory/" + sender + "/mytoken", "bad", true},
		{"non-factory denom", sender, "uheya", other, true},
		{"empty denom", sender, "", other, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgChangeAdmin(tt.sender, tt.denom, tt.newAdmin)
			err := msg.ValidateBasic()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMsgForceTransfer_ValidateBasic(t *testing.T) {
	sender := validSenderAddr(t)
	other := otherSenderAddr(t)
	denom := "factory/" + sender + "/mytoken"

	tests := []struct {
		name        string
		sender      string
		amount      string
		destAddr    string
		fromAddress string
		expectErr   bool
	}{
		{"valid", sender, "1000" + denom, other, "", false},
		{"invalid sender", "bad", "1000" + denom, other, "", true},
		{"invalid dest", sender, "1000" + denom, "bad", "", true},
		{"non-factory denom", sender, "1000uheya", other, "", true},
		{"zero amount", sender, "0" + denom, other, "", true},
		{"invalid from address", sender, "1000" + denom, other, "bad", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgForceTransferFull(tt.sender, tt.amount, tt.destAddr, tt.fromAddress)
			err := msg.ValidateBasic()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewDenom(t *testing.T) {
	denom := types.NewDenom("heya1abc", "mytoken")
	expected := "factory/heya1abc/mytoken"
	if denom != expected {
		t.Errorf("expected %s, got %s", expected, denom)
	}
}

func TestGetSigners(t *testing.T) {
	sender := validSenderAddr(t)
	addr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		t.Fatal(err)
	}

	msg := types.NewMsgCreateDenom(sender, "test")
	signers := msg.GetSigners()
	if len(signers) != 1 || !signers[0].Equals(addr) {
		t.Errorf("unexpected signers")
	}
}
