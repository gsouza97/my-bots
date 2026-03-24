package composer

import "github.com/gsouza97/my-bots/internal/infrastructure/providers"

type ProvidersComposer struct {
	BinancePriceProvider  *providers.BinancePriceProvider
	RevertFeeProvider     *providers.RevertFeeProvider
	FearAndGreedProvider  *providers.AlternativeFearAndGreedProvider
	AltcoinSeasonProvider *providers.CmcAltcoinSeasonProvider
	HomologacionProvider  *providers.HomologacionProvider
}

func NewProvidersComposer() *ProvidersComposer {
	return &ProvidersComposer{
		BinancePriceProvider:  providers.NewBinancePriceProvider(),
		RevertFeeProvider:     providers.NewRevertFeeProvider(),
		FearAndGreedProvider:  providers.NewAlternativeFearAndGreedProvider(),
		AltcoinSeasonProvider: providers.NewCmcAltcoinSeasonProvider(),
		HomologacionProvider:  providers.NewHomologacionProvider(),
	}
}
