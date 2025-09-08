package main

func handleSelection(ordered []FundInfo, unrateTrigger bool) error {
	selected, err := selector(ordered, unrateTrigger)
	if err != nil {
		return err
	}

	print(selected)

	return nil
}

// selector selects funds based on whether the unrateTrigger is true or false
// If unrateTrigger is false, it selects the top 3 funds regardless of their SMA status
// If unrateTrigger is true, it selects the top 3 funds whose NAV is above their SMA
func selector(ordered []FundInfo, unrateTrigger bool) ([]FundInfo, error) {
	var selected []FundInfo

	// If unrate is above its MA12 any selected fund must have its NAV above its SMA
	if !unrateTrigger { // grab top 3 regardless of SMA
		selected = append(selected, ordered[:3]...)
	}
	if unrateTrigger { // grab top 3 that have NAV above SMA
		limit := 3
		var count int

		for _, f := range ordered {
			if count >= limit {
				break
			}
			if f.SmaDistance > 0 {
				selected = append(selected, f)
				count++
			}
		}
	}

	return selected, nil
}
