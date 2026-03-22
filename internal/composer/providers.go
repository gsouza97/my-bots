package composer

import "github.com/gsouza97/my-bots/internal/adapter/provider"

type ProvidersComposer struct {
	BinancePriceProvider  *provider.BinancePriceProvider
	RevertFeeProvider     *provider.RevertFeeProvider
	FearAndGreedProvider  *provider.AlternativeFearAndGreedProvider
	AltcoinSeasonProvider *provider.CmcAltcoinSeasonProvider
	HomologacionProvider  *provider.HomologacionProvider
}

func NewProvidersComposer() *ProvidersComposer {
	return &ProvidersComposer{
		BinancePriceProvider:  provider.NewBinancePriceProvider(),
		RevertFeeProvider:     provider.NewRevertFeeProvider(),
		FearAndGreedProvider:  provider.NewAlternativeFearAndGreedProvider(),
		AltcoinSeasonProvider: provider.NewCmcAltcoinSeasonProvider(),
		HomologacionProvider:  provider.NewHomologacionProvider(),
	}
}
