package main

// TODO mock server so that we dont need to reach out to the internet
// Likely easiest by feeding the retrieval function an URL (or a function that returns an URL)
/*
import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestRetrieveData(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()

	fund := Fund{
		AzaID: 1534742,
	}

	fundData, err := retrieveFundData(t.Context(), fund.AzaID)
	assert.NoError(t, err)
	assert.NotZero(t, fundData.From)
}
*/
