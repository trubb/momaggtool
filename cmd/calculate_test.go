package main

import (
	"testing"
)

// github.com/stretchr/testify/require

func TestCalculate(t *testing.T) {
	t.Parallel()

	// TODO use realer data as well

	tests := []struct {
		name     string
		period   int
		data     []AzaOhlc
		expected float64
	}{
		{
			name:   "simple case",
			period: 3,
			data: []AzaOhlc{
				{Close: 100.0},
				{Close: 200.0},
				{Close: 300.0},
			},
			expected: 200.0,
		},
		{
			name:   "longer period",
			period: 5,
			data: []AzaOhlc{
				{Close: 100.0},
				{Close: 200.0},
				{Close: 300.0},
				{Close: 400.0},
				{Close: 500.0},
			},
			expected: 300.0,
		},
		{
			name:   "not enough data",
			period: 5,
			data: []AzaOhlc{
				{Close: 100.0},
				{Close: 200.0},
				{Close: 300.0},
			},
			expected: 0.0,
		},
		{
			name:   "no funny business",
			period: 10,
			data: []AzaOhlc{
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
				{Close: 100.0},
			},
			expected: 100.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := simpleMovingAverage(tt.period, tt.data, tt.name)
			if len(tt.data) < tt.period {
				if err == nil {
					t.Errorf("expected error for not enough data, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if f != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, f)
			}
		})
	}

}
