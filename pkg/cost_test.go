package pkg

import (
	"math"
	"reflect"
	"testing"
)

func TestNewCostAnalysis(t *testing.T) {
	tests := []struct {
		name          string
		priceLocation string
		expected      *CostAnalysis
		expectedError bool
	}{
		{
			name:          "blank",
			priceLocation: "",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "nonexistent price path",
			priceLocation: "testdata/i_dont_exist.json",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "malformed json local",
			priceLocation: "testdata/im_not_json.json",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "valid local pricing",
			priceLocation: "testdata/valid_pricing.json",
			expected: &CostAnalysis{
				priceSheetPath: "testdata/valid_pricing.json",
				pricing: Pricing{
					"us-west1-a": {
						"us-west1-b": 0.01,
					},
					"us-west1-b": {
						"us-west1-a": 0.01,
					},
				},
			},
		},
		{
			name:          "nonexistent url",
			priceLocation: "https://goo.bar/elmo.json",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid remote pricing",
			priceLocation: "https://raw.githubusercontent.com/tetratelabs/istio-cost-analyzer/cost-unit-tests/pkg/testdata/im_not_json.json",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "valid remote pricing",
			priceLocation: "https://raw.githubusercontent.com/tetratelabs/istio-cost-analyzer/cost-unit-tests/pkg/testdata/valid_pricing.json",
			expected: &CostAnalysis{
				priceSheetPath: "https://raw.githubusercontent.com/tetratelabs/istio-cost-analyzer/cost-unit-tests/pkg/testdata/valid_pricing.json",
				pricing: Pricing{
					"us-west1-a": {
						"us-west1-b": 0.01,
					},
					"us-west1-b": {
						"us-west1-a": 0.01,
					},
				},
			},
			expectedError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ca, err := NewCostAnalysis(tt.priceLocation); ((err != nil) != tt.expectedError) || !reflect.DeepEqual(ca, tt.expected) {
				t.Errorf("expected error existence: %v => (%v), expected CostAnalysis object %v => (%v)", tt.expectedError, err != nil, tt.expected, ca)
			}
		})
	}
}

func TestCostAnalysis_CalculateEgress(t *testing.T) {
	tests := []struct {
		name           string
		callsWithPrice []*Call
		expectedTotal  float64
		expectedError  bool
	}{
		{
			name:           "empty calls",
			callsWithPrice: make([]*Call, 0),
			expectedTotal:  0.00,
			expectedError:  false,
		},
		{
			name: "non-existent regions",
			callsWithPrice: []*Call{
				{
					From:     "us-east1-b",
					To:       "us-west1-b",
					CallSize: uint64(math.Pow(10, 6)),
				},
			},
			expectedTotal: 0.00,
			expectedError: false,
		},
		{
			name: "legit prices",
			callsWithPrice: []*Call{
				{
					From:     "us-west1-b",
					To:       "us-east1-b",
					CallSize: uint64(math.Pow(10, 6)),
					CallCost: 0.9,
				},
				{
					From:     "us-west1-b",
					To:       "us-west1-c",
					CallSize: uint64(math.Pow(10, 6)),
					CallCost: 0.5,
				},
			},
			expectedTotal: 1.4,
			expectedError: false,
		},
	}
	ca := &CostAnalysis{
		pricing: Pricing{
			"us-west1-b": {
				"us-west1-c": 0.5,
				"us-east1-b": 0.9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strippedCalls := make([]*Call, 0)
			for _, v := range tt.callsWithPrice {
				stripped := *v
				stripped.CallCost = 0.00
				strippedCalls = append(strippedCalls, &stripped)
			}
			total, err := ca.CalculateEgress(strippedCalls)
			if total != tt.expectedTotal || (err != nil) != tt.expectedError ||
				!reflect.DeepEqual(tt.callsWithPrice, strippedCalls) {
				t.Errorf("expected err (%v)=>%v, expected total (%v)=>%v, expected pricing (%v)=>%v\n",
					tt.expectedError, err != nil, tt.expectedTotal, total, tt.callsWithPrice, strippedCalls)
			}
		})
	}
}
