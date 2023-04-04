package keeper

import (
	"bytes"
	"encoding/hex"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"nucleus/x/htlc/types"
)

// CreateHTLC creates an HTLC
func (k Keeper) CreateHTLC(
	ctx sdk.Context,
	sender sdk.AccAddress,
	to sdk.AccAddress,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
	timestamp uint64,
	timeLock uint64,
) (
	id tmbytes.HexBytes,
	err error,
) {
	id = types.GetID(sender, to, amount, hashLock)

	// check if the HTLC already exists
	if k.HasHTLC(ctx, id) {
		return id, sdkerrors.Wrap(types.ErrHTLCExists, id.String())
	}

	expirationHeight := uint64(ctx.BlockHeight()) + timeLock

	if err = k.createHTLC(ctx, sender, amount); err != nil {
		return id, err
	}

	htlc := types.NewHTLC(
		id, sender, to, amount, hashLock,
		nil, timestamp, expirationHeight,
		types.Open, 0,
	)

	// set the HTLC
	k.SetHTLC(ctx, htlc, id)

	// add to the expiration queue
	k.AddHTLCToExpiredQueue(ctx, htlc.ExpirationHeight, id)

	return id, nil
}

func (k Keeper) createHTLC(
	ctx sdk.Context,
	sender sdk.AccAddress,
	amount sdk.Coins,
) error {
	// transfer the specified tokens to the HTLC module account
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount)
}

// ClaimHTLC claims the specified HTLC with the given secret
func (k Keeper) ClaimHTLC(
	ctx sdk.Context,
	id tmbytes.HexBytes,
	secret tmbytes.HexBytes,
) (
	string,
	error,
) {
	// query the HTLC
	htlc, found := k.GetHTLC(ctx, id)
	if !found {
		return "", sdkerrors.Wrap(types.ErrUnknownHTLC, id.String())
	}

	// check if the HTLC is open
	if htlc.State != types.Open {
		return "", sdkerrors.Wrap(types.ErrHTLCNotOpen, id.String())
	}

	hashLock, _ := hex.DecodeString(htlc.HashLock)

	// check if the secret matches with the hash lock
	if !bytes.Equal(types.GetHashLock(secret, htlc.Timestamp), hashLock) {
		return "", sdkerrors.Wrap(types.ErrInvalidSecret, secret.String())
	}

	to, err := sdk.AccAddressFromBech32(htlc.To)
	if err != nil {
		return "", err
	}

	if err := k.claimHTLC(ctx, htlc.Amount, to); err != nil {
		return "", err
	}

	// update the secret and state of the HTLC
	htlc.Secret = secret.String()
	htlc.State = types.Completed
	htlc.ClosedBlock = uint64(ctx.BlockHeight())
	k.SetHTLC(ctx, htlc, id)

	// delete from the expiration queue
	k.DeleteHTLCFromExpiredQueue(ctx, htlc.ExpirationHeight, id)

	return htlc.HashLock, nil
}

func (k Keeper) claimHTLC(ctx sdk.Context, amount sdk.Coins, to sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, amount)
}

// RefundHTLC refunds the specified HTLC
func (k Keeper) RefundHTLC(ctx sdk.Context, h types.HTLC, id tmbytes.HexBytes) error {
	sender, err := sdk.AccAddressFromBech32(h.Sender)
	if err != nil {
		return err
	}

	if err := k.refundHTLC(ctx, sender, h.Amount); err != nil {
		return err
	}

	// update the state of the HTLC
	h.State = types.Refunded
	h.ClosedBlock = uint64(ctx.BlockHeight())
	k.SetHTLC(ctx, h, id)

	return nil
}

func (k Keeper) refundHTLC(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, amount)
}

// HasHTLC checks if the given HTLC exists
func (k Keeper) HasHTLC(ctx sdk.Context, id tmbytes.HexBytes) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetHTLCKey(id))
}

// SetHTLC sets the given HTLC
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, id tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&htlc)
	store.Set(types.GetHTLCKey(id), bz)
}

// GetHTLC retrieves the specified HTLC
func (k Keeper) GetHTLC(ctx sdk.Context, id tmbytes.HexBytes) (htlc types.HTLC, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetHTLCKey(id))
	if bz == nil {
		return htlc, false
	}
	k.cdc.MustUnmarshal(bz, &htlc)
	return htlc, true
}

// AddHTLCToExpiredQueue adds the specified HTLC to the expiration queue
func (k Keeper) AddHTLCToExpiredQueue(ctx sdk.Context, expirationHeight uint64, id tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetHTLCExpiredQueueKey(expirationHeight, id), []byte{})
}

// DeleteHTLCFromExpiredQueue removes the specified HTLC from the expiration queue
func (k Keeper) DeleteHTLCFromExpiredQueue(ctx sdk.Context, expirationHeight uint64, id tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetHTLCExpiredQueueKey(expirationHeight, id))
}

// IterateHTLCs iterates through the HTLCs
func (k Keeper) IterateHTLCs(
	ctx sdk.Context,
	op func(id tmbytes.HexBytes, h types.HTLC) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.HTLCKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := tmbytes.HexBytes(iterator.Key()[1:])

		var htlc types.HTLC
		k.cdc.MustUnmarshal(iterator.Value(), &htlc)

		if stop := op(id, htlc); stop {
			break
		}
	}
}

// IterateHTLCExpiredQueueByHeight iterates through the HTLC expiration queue by the specified height
func (k Keeper) IterateHTLCExpiredQueueByHeight(
	ctx sdk.Context, height uint64,
	op func(id tmbytes.HexBytes, h types.HTLC) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetHTLCExpiredQueueSubspace(height))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := tmbytes.HexBytes(iterator.Key()[9:])
		htlc, _ := k.GetHTLC(ctx, id)

		if stop := op(id, htlc); stop {
			break
		}
	}
}
