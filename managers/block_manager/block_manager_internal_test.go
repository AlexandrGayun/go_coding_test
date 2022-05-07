package block_manager

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseJsonResponse(t *testing.T) {
	JSONExampleTests := []struct {
		jsonBytes       []byte
		responseStrings []string
	}{
		{[]byte(`{"SomeField": "SomeValue", "result": { "SomeField": "SomeValue", "transactions": [{ "SomeField": "SomeValue", "value": "0x137e9165c0e3000"}, {"value": "0x431a5d8aa58d000"}]}}`), []string{"0x137e9165c0e3000", "0x431a5d8aa58d000"}},
		{[]byte(`{"result": { "transactions": [{ "value": "0x137e9165c0e3000"}, {"value": "0x431a5d8aa58d000"}]}}`), []string{"0x137e9165c0e3000", "0x431a5d8aa58d000"}},
		{[]byte(`{"transactions": [{ "value": "0x137e9165c0e3000"}, {"value": "0x431a5d8aa58d000"}]}`), []string{}},
		{[]byte(`{"error": "errorMessage"}`), []string{}},
	}
	for _, tt := range JSONExampleTests {
		t.Run("parse response body to extract values ", func(t *testing.T) {
			strArr, err := parseJSONResponse(tt.jsonBytes)
			if err != nil {
				t.Error("got fastjson parser error", err)
			}
			if !reflect.DeepEqual(strArr, tt.responseStrings) {
				t.Errorf("wrong json interpretation. got %q, want %q", strArr, tt.responseStrings)
			}
		})
	}
}

func TestGroupTransactions(t *testing.T) {
	dataToGroup := []struct {
		inputStrings  []string
		resultTCount  int
		resultTAmount float64
	}{
		{[]string{"0x137e9165c0e3000", "0x431a5d8aa58d000"}, 2, 0.390000},
		{[]string{"0x137e9165c0e3000"}, 1, 0.087795},
		{[]string{}, 0, 0},
	}
	for _, tt := range dataToGroup {
		t.Run("should count correctly", func(t *testing.T) {
			resTCount, resTAmount, err := groupTransactions(tt.inputStrings)
			if err != nil {
				t.Error("got unexpected error", err)
			}
			if resTCount != tt.resultTCount ||
				fmt.Sprintf("%.6f", resTAmount) != fmt.Sprintf("%.6f", tt.resultTAmount) {
				t.Errorf(`got unexpected results. Total amount should be %.6f, got %.6f
								Total transactions should be %v, got %v`,
					tt.resultTAmount,
					resTAmount,
					tt.resultTCount,
					resTCount)
			}
		})
	}
}
