package keeper_test

import (
	"testing"

	testkeeper "github.com/komodoplatform/nucleus/testutil/keeper"
	"github.com/komodoplatform/nucleus/x/nucleus/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.NucleusKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
