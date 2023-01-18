package nucleus_test

import (
	"testing"

	keepertest "nucleus/testutil/keeper"
	"nucleus/testutil/nullify"
	"nucleus/x/nucleus"
	"nucleus/x/nucleus/types"

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
