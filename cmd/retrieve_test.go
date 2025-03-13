package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// github.com/stretchr/testify/require

func TestRetrieveData(t *testing.T) {
	fund := Fund{
		AzaID: 1534742,
	}

	fundData, err := retrieveFundData(fund)
	assert.NoError(t, err)
	assert.NotZero(t, fundData.From)
	// assert.???
}

// TODO mock server so that we dont need to reach out to the internet
