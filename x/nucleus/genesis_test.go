package nucleus_test

import (
	"testing"

	keepertest "github.com/komodoplatform/nucleus/testutil/keeper"
	"github.com/komodoplatform/nucleus/testutil/nullify"
	"github.com/komodoplatform/nucleus/x/nucleus"
	"github.com/komodoplatform/nucleus/x/nucleus/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keepertest.NucleusKeeper(t)
	nucleus.InitGenesis(ctx, *k, genesisState)
	got := nucleus.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
