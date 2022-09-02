package entities

type (
	BuyOrder struct {
		UnitBought int64
	}

	BuySignal struct {
		ShouldBuy bool
		Order     BuyOrder
	}
)
