package srcscan

import (
	"encoding/json"
	"github.com/kr/pretty"
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	type unmarshalJSONTest struct {
		dir string
	}
	tests := []unmarshalJSONTest{
		{"testdata"},
	}
	for _, test := range tests {
		units, err := Default.Scan(test.dir)
		if err != nil {
			t.Errorf("scan error: %s", err)
			continue
		}
		for _, unit := range units {
			data, err := json.Marshal(unit)
			if err != nil {
				t.Errorf("marshal error: %s", err)
				continue
			}
			unit2, err := UnmarshalJSON(data, UnitType(unit))
			if err != nil {
				t.Errorf("UnmarshalJSON error: %s", err)
				continue
			}
			if !reflect.DeepEqual(unit, unit2) {
				t.Errorf("unit != unit2:\n%+v\n%+v\n%v", unit, unit2, strings.Join(pretty.Diff(unit, unit2), "\n"))
			}
		}
	}
}

func TestMarshalableUnit(t *testing.T) {
	type unmarshalJSONTest struct {
		dir string
	}
	tests := []unmarshalJSONTest{
		{"testdata"},
	}
	for _, test := range tests {
		units, err := Default.Scan(test.dir)
		if err != nil {
			t.Errorf("scan error: %s", err)
			continue
		}
		for _, unit := range units {
			mu := &MarshalableUnit{unit}
			data, err := json.Marshal(mu)
			if err != nil {
				t.Errorf("marshal error: %s", err)
				continue
			}
			var mu2 *MarshalableUnit
			err = json.Unmarshal(data, &mu2)
			if err != nil {
				t.Errorf("Unmarshal error: %s", err)
				continue
			}
			if !reflect.DeepEqual(mu, mu2) {
				t.Errorf("mu != mu2:\n%+v\n%+v\n%v", mu, mu2, strings.Join(pretty.Diff(mu, mu2), "\n"))
			}
		}
	}
}
