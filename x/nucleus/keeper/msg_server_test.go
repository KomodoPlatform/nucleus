package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/komodoplatform/nucleus/testutil/keeper"
	"github.com/komodoplatform/nucleus/x/nucleus/keeper"
	"github.com/komodoplatform/nucleus/x/nucleus/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO
// nolint
func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NucleusKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
