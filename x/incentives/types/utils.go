package types

import "sort"

// IncentivesToBytes converts a slice of Incentives to a slice of bytes.
func IncentivesToBytes(incentives ...Incentive) ([][]byte, error) {
	incentiveBytes := make([][]byte, len(incentives))
	for i, incentive := range incentives {
		bz, err := incentive.Marshal()
		if err != nil {
			return nil, err
		}

		incentiveBytes[i] = bz
	}

	return incentiveBytes, nil
}

// SortIncentivesStrategiesMap sorts a map of IncentivesByType by their type.
func SortIncentivesStrategiesMap(incentiveStrategies map[Incentive]Strategy) []Incentive {
	// Get all incentive types and sort them by name.
	incentiveTypes := make([]Incentive, len(incentiveStrategies))
	i := 0
	for incentive := range incentiveStrategies {
		incentiveTypes[i] = incentive
		i++
	}

	// Sort the incentive types by name.
	return SortIncentivesByType(incentiveTypes)
}

// SortIncentivesByType sorts a slice of Incentives by their type.
func SortIncentivesByType(incentives []Incentive) []Incentive {
	// Sort the incentive types by name.
	sort.Slice(incentives, func(i, j int) bool {
		return incentives[i].Type() < incentives[j].Type()
	})

	return incentives
}
