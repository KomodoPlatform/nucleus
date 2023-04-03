package htlc_test

import (
	"math/rand"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"nucleus/x/htlc/types"
)

var (
	DenomMap                    = map[int]string{0: "htltbnb", 1: "htltinc"}
	MinTimeLock          uint64 = 220
	MaxTimeLock          uint64 = 270
	TestDeputy                  = sdk.AccAddress(crypto.AddressHash([]byte("TestDeputy")))
)

// TODO
// nolint
func c(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}

// TODO
// nolint
func cs(coins ...sdk.Coin) sdk.Coins {
	return sdk.NewCoins(coins...)
}

// TODO
// nolint
func ts(minOffset int) uint64 {
	return uint64(time.Now().Add(time.Duration(minOffset) * time.Minute).Unix())
}

// GeneratePrivKeyAddressPairsFromRand generates (deterministically) a total of n secp256k1 private keys and addresses.
func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	r := rand.New(rand.NewSource(12345)) // make the generation deterministic
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		if _, err := r.Read(secret); err != nil {
			panic("Could not read randomness")
		}
		keys[i] = secp256k1.GenPrivKeySecp256k1(secret)
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

func GenerateRandomSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

// TODO
// nolint
func loadSwapAndSupply(addr sdk.AccAddress, index int) (types.HTLC, types.AssetSupply) {
	coin := c(DenomMap[index], 50000)
	expireOffset := MinTimeLock
	timestamp := ts(index)
	amount := cs(coin)
	randomSecret, _ := GenerateRandomSecret()
	randomHashLock := types.GetHashLock(randomSecret, timestamp)
	id := types.GetID(addr, addr, amount, randomHashLock)
	htlc := types.NewHTLC(
		id,
		addr,
		addr,
		amount,
		randomHashLock,
		[]byte{},
		timestamp,
		expireOffset,
		types.Open,
		1,
	)

	supply := types.NewAssetSupply(
		coin,
		c(coin.Denom, 0),
		c(coin.Denom, 0),
		c(coin.Denom, 0),
		time.Duration(0),
	)

	return htlc, supply
}
