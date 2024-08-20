package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func (s *KeeperTestSuite) TestInitGenesis() {
	badSLA := slatypes.PriceFeedSLA{}

	sla1 := slatypes.NewPriceFeedSLA(
		"id",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	sla2 := slatypes.NewPriceFeedSLA(
		"id2",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	sla3 := slatypes.NewPriceFeedSLA(
		"id3",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	cp1 := slinkytypes.NewCurrencyPair("BTC", "USD")

	consAddress1 := sdk.ConsAddress("consAddress1")
	consAddress2 := sdk.ConsAddress("consAddress2")

	priceFeed1, err := slatypes.NewPriceFeed(
		10,
		consAddress1,
		cp1,
		"id1",
	)
	s.Require().NoError(err)
	priceFeed2, _ := slatypes.NewPriceFeed(
		10,
		consAddress2,
		cp1,
		"id1",
	)
	s.Require().NoError(err)

	s.Run("bad genesis state should panic", func() {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{badSLA}, nil, slatypes.DefaultParams())
		s.Require().Panics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})
	})

	s.Run("good genesis state should not panic", func() {
		gs := slatypes.NewGenesisState([]slatypes.PriceFeedSLA{sla1, sla2, sla3}, nil, slatypes.DefaultParams())

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(3, len(slas))
		s.Require().Equal(sla1, slas[0])
		s.Require().Equal(sla2, slas[1])
		s.Require().Equal(sla3, slas[2])

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(slatypes.DefaultParams(), params)

		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(0, len(cps))
	})

	s.Run("default genesis state", func() {
		gs := slatypes.NewDefaultGenesisState()

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(0, len(slas))

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(slatypes.DefaultParams(), params)

		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(0, len(cps))
	})

	s.Run("can init with slas and price feeds", func() {
		gs := slatypes.NewDefaultGenesisState()
		gs.PriceFeeds = []slatypes.PriceFeed{priceFeed1, priceFeed2}

		sla1 := slatypes.NewPriceFeedSLA(
			priceFeed1.ID,
			10,
			math.LegacyMustNewDecFromStr("1.0"),
			math.LegacyMustNewDecFromStr("1.0"),
			5,
			5,
		)

		gs.SLAs = []slatypes.PriceFeedSLA{sla1}

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(slas))
		s.Require().Equal(sla1, slas[0])

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, sla1.ID)
		s.Require().NoError(err)
		s.Require().Equal(2, len(feeds))

		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))
		s.Require().Contains(cps, priceFeed1.CurrencyPair)
		s.Require().Contains(cps, priceFeed2.CurrencyPair)
	})

	s.Run("can init with slas and price feeds and updated params", func() {
		gs := slatypes.NewDefaultGenesisState()
		gs.PriceFeeds = []slatypes.PriceFeed{priceFeed1, priceFeed2}

		sla1 := slatypes.NewPriceFeedSLA(
			priceFeed1.ID,
			10,
			math.LegacyMustNewDecFromStr("1.0"),
			math.LegacyMustNewDecFromStr("1.0"),
			5,
			5,
		)

		gs.SLAs = []slatypes.PriceFeedSLA{sla1}

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})

		slas, err := s.keeper.GetSLAs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(slas))
		s.Require().Equal(sla1, slas[0])

		feeds, err := s.keeper.GetAllPriceFeeds(s.ctx, sla1.ID)
		s.Require().NoError(err)
		s.Require().Equal(2, len(feeds))

		gs.Params = slatypes.Params{
			Enabled: false,
		}

		s.Require().NotPanics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})

		params, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(slatypes.Params{Enabled: false}, params)

		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(cps))
		s.Require().Contains(cps, priceFeed1.CurrencyPair)
		s.Require().Contains(cps, priceFeed2.CurrencyPair)
	})

	s.Run("bad price feed", func() {
		gs := slatypes.NewDefaultGenesisState()
		gs.PriceFeeds = []slatypes.PriceFeed{priceFeed1, priceFeed2}

		sla1 := slatypes.NewPriceFeedSLA(
			"",
			0,
			math.LegacyMustNewDecFromStr("1.0"),
			math.LegacyMustNewDecFromStr("1.0"),
			5,
			5,
		)

		gs.SLAs = []slatypes.PriceFeedSLA{sla1}

		s.Require().Panics(func() {
			s.keeper.InitGenesis(s.ctx, *gs)
		})
	})
}

func (s *KeeperTestSuite) TestExportGenesis() {
	sla1 := slatypes.NewPriceFeedSLA(
		"id1",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	sla2 := slatypes.NewPriceFeedSLA(
		"id2",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	sla3 := slatypes.NewPriceFeedSLA(
		"id3",
		10,
		math.LegacyMustNewDecFromStr("1.0"),
		math.LegacyMustNewDecFromStr("1.0"),
		5,
		5,
	)

	cp1 := slinkytypes.NewCurrencyPair("btc", "usd")

	consAddress1 := sdk.ConsAddress("consAddress1")
	consAddress2 := sdk.ConsAddress("consAddress2")

	priceFeed1, err := slatypes.NewPriceFeed(
		10,
		consAddress1,
		cp1,
		"id1",
	)
	s.Require().NoError(err)
	priceFeed2, _ := slatypes.NewPriceFeed(
		10,
		consAddress2,
		cp1,
		"id1",
	)
	s.Require().NoError(err)

	s.Run("can export an empty genesis state", func() {
		gs := s.keeper.ExportGenesis(s.ctx)
		defaultGS := slatypes.NewDefaultGenesisState()

		s.Require().Equal(defaultGS.Params, gs.Params)
		s.Require().Equal(len(defaultGS.SLAs), len(gs.SLAs))
		s.Require().Equal(len(defaultGS.PriceFeeds), len(gs.PriceFeeds))
	})

	s.Run("can export with slas", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla1, sla2, sla3})
		s.Require().NoError(err)

		gs := s.keeper.ExportGenesis(s.ctx)
		defaultGS := slatypes.NewDefaultGenesisState()
		defaultGS.SLAs = []slatypes.PriceFeedSLA{sla1, sla2, sla3}

		s.Require().Equal(defaultGS.Params, gs.Params)
		s.Require().Equal(len(defaultGS.SLAs), len(gs.SLAs))
		s.Require().Equal(len(defaultGS.PriceFeeds), len(gs.PriceFeeds))

		for i := range gs.SLAs {
			s.Require().Equal(defaultGS.SLAs[i], gs.SLAs[i])
		}
	})

	s.Run("can export with slas and updated params", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla1, sla2, sla3})
		s.Require().NoError(err)

		params := slatypes.Params{Enabled: false}
		err = s.keeper.SetParams(s.ctx, params)
		s.Require().NoError(err)

		gs := s.keeper.ExportGenesis(s.ctx)
		defaultGs := slatypes.NewDefaultGenesisState()
		defaultGs.SLAs = []slatypes.PriceFeedSLA{sla1, sla2, sla3}
		defaultGs.Params = params

		s.Require().Equal(defaultGs.Params, gs.Params)
		s.Require().Equal(len(defaultGs.SLAs), len(gs.SLAs))
		s.Require().Equal(len(defaultGs.PriceFeeds), len(gs.PriceFeeds))

		for i := range gs.SLAs {
			s.Require().Equal(defaultGs.SLAs[i], gs.SLAs[i])
		}
	})

	s.Run("can export with slas and price feeds", func() {
		err := s.keeper.AddSLAs(s.ctx, []slatypes.PriceFeedSLA{sla1, sla2, sla3})
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed1)
		s.Require().NoError(err)

		err = s.keeper.SetPriceFeed(s.ctx, priceFeed2)
		s.Require().NoError(err)

		gs := s.keeper.ExportGenesis(s.ctx)
		defaultGs := slatypes.NewDefaultGenesisState()
		defaultGs.SLAs = []slatypes.PriceFeedSLA{sla1, sla2, sla3}
		defaultGs.PriceFeeds = []slatypes.PriceFeed{priceFeed1, priceFeed2}

		s.Require().Equal(defaultGs.Params, gs.Params)
		s.Require().Equal(len(defaultGs.SLAs), len(gs.SLAs))
		s.Require().Equal(len(defaultGs.PriceFeeds), len(gs.PriceFeeds))

		for i := range gs.SLAs {
			s.Require().Equal(defaultGs.SLAs[i], gs.SLAs[i])
		}

		for i := range gs.PriceFeeds {
			s.Require().Equal(defaultGs.PriceFeeds[i], gs.PriceFeeds[i])
		}
	})
}
