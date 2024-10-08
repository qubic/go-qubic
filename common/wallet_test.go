//go:build (linux || darwin) && amd64
// +build linux darwin
// +build amd64

package common

import (
	"encoding/hex"
	"github.com/cloudflare/circl/ecc/fourq"
	fourq2 "github.com/cloudflare/fourq"
	"github.com/google/go-cmp/cmp"
	"testing"
)

const testSeed = "lmujdbtiimznvyamoyjumfeiglauqfzsznisydmqrdyuwqydjpjixat"

func TestGetSubseed(t *testing.T) {
	expectedSubseedBytes := [32]byte{244, 124, 158, 118, 114, 22, 0, 127, 168, 254, 156, 41, 77, 119, 59, 224, 162, 60, 9, 187, 43, 141, 1, 189, 213, 224, 195, 24, 35, 144, 238, 58}

	got, err := GetSubSeed(testSeed)
	if err != nil {
		t.Fatalf("Got err when getting subseed. err: %s", err.Error())
	}

	if cmp.Diff(got, expectedSubseedBytes) != "" {
		t.Fatalf("Mismatched return value. Expected: %s, got: %s", hex.EncodeToString(expectedSubseedBytes[:]), hex.EncodeToString(got[:]))
	}
}

func TestGetPrivateKey(t *testing.T) {
	expectedPrivKey := [32]byte{255, 152, 128, 102, 167, 172, 117, 67, 207, 98, 121, 87, 47, 195, 144, 191, 211, 225, 145, 187, 93, 83, 248, 238, 217, 120, 166, 88, 206, 146, 124, 225}
	got, err := getPrivateKey(testSeed)
	if err != nil {
		t.Fatalf("Got err when getting priv key. err: %s", err.Error())
	}

	if cmp.Diff(got, expectedPrivKey) != "" {
		t.Fatalf("Mismatched return value. Expected: %s, got: %s", hex.EncodeToString(expectedPrivKey[:]), hex.EncodeToString(got[:]))
	}
}

func TestGetPublicKey(t *testing.T) {
	privKey := [32]byte{255, 152, 128, 102, 167, 172, 117, 67, 207, 98, 121, 87, 47, 195, 144, 191, 211, 225, 145, 187, 93, 83, 248, 238, 217, 120, 166, 88, 206, 146, 124, 225}
	expectedPubKey := [32]byte{230, 252, 58, 173, 75, 89, 77, 130, 191, 49, 3, 161, 16, 22, 216, 13, 232, 131, 222, 135, 59, 206, 196, 142, 144, 57, 98, 134, 80, 59, 38, 19}
	got, err := getPublicKey(privKey)
	if err != nil {
		t.Fatalf("Got err when getting pub key. err: %s", err.Error())
	}

	if cmp.Diff(got, expectedPubKey) != "" {
		t.Fatalf("Mismatched return value. Expected: %s, got: %s", hex.EncodeToString(expectedPubKey[:]), hex.EncodeToString(got[:]))
	}
}

func TestCreateWallet(t *testing.T) {
	expected := Wallet{
		PubKey:   [32]byte{230, 252, 58, 173, 75, 89, 77, 130, 191, 49, 3, 161, 16, 22, 216, 13, 232, 131, 222, 135, 59, 206, 196, 142, 144, 57, 98, 134, 80, 59, 38, 19},
		PrivKey:  [32]byte{255, 152, 128, 102, 167, 172, 117, 67, 207, 98, 121, 87, 47, 195, 144, 191, 211, 225, 145, 187, 93, 83, 248, 238, 217, 120, 166, 88, 206, 146, 124, 225},
		Identity: "QJRRSSKMJRDKUDTYVNYGAMQPULKAMILQQYOWBEXUDEUWQUMNGDHQYLOAJMEB",
	}

	got, err := NewWallet(testSeed)
	if err != nil {
		t.Fatalf("Got err when creating wallet. err: %s", err.Error())
	}

	if diff := cmp.Diff(got, expected); diff != "" {
		t.Fatalf("Mismatched return value. Diff: %s", diff)
	}
}

func TestDecodePubKey(t *testing.T) {
	privKey, err := getPrivateKey(testSeed)
	if err != nil {
		t.Fatalf("err creating privkey")
	}

	var pubKeyPoint fourq.Point
	pubKeyPoint.ScalarBaseMult(&privKey)

	pubKey, err := encode(pubKeyPoint)
	if err != nil {
		t.Fatalf("err encoding pubkey")
	}

	ok := fourq2.IsOnCurve(pubKey)
	if !ok {
		t.Fatalf("pubkey not on curve")
	}
}
