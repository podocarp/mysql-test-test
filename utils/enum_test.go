package utils_test

import (
	"slices"
	"testing"

	"github.com/podocarp/mysql-test-test/utils"
)

func assertBitsetEq(t *testing.T, expected, actual utils.CountryBitset) {
	t.Helper()
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("Expected %v, obtained %v", expected, actual)
		}
	}
}

func assertCountriesEq(t *testing.T, expected, actual utils.Countries) {
	t.Helper()
	fail := false
	defer func() {
		t.Helper()
		if fail {
			t.Fatalf("Expected %v, obtained %v", expected, actual)
		}
	}()

	if len(expected) != len(actual) {
		fail = true
		return
	}

	slices.Sort(expected)
	slices.Sort(actual)
	for i := range expected {
		if expected[i] != actual[i] {
			fail = true
			return
		}
	}
}

func TestToBitset(t *testing.T) {
	countries := utils.Countries{
		utils.COUNTRY_AD,
		utils.COUNTRY_FR,
		utils.COUNTRY_AG,
		utils.COUNTRY_NU,
	}
	res := countries.ToBitset()
	ans := utils.CountryBitset{
		1<<5 | 1<<9, // AD and AG
		1 << 12,     // FR
		1 << 33,     // NU
		0,
	}
	assertBitsetEq(t, ans, res)

	resCountries := res.ToCountries()
	assertCountriesEq(t, countries, resCountries)
}

func FuzzToBitset(f *testing.F) {
	f.Add(uint64(1), uint64(1), uint64(1), uint64(1))
	f.Fuzz(func(t *testing.T, a, b, c, d uint64) {
		bitset := utils.CountryBitset{a, b, c, d}
		countries := bitset.ToCountries()
		res := countries.ToBitset()
		assertBitsetEq(t, bitset, res)
	})
}

func FuzzValuerScanner(f *testing.F) {
	f.Add(uint64(1), uint64(1), uint64(1), uint64(1))
	f.Fuzz(func(t *testing.T, a, b, c, d uint64) {
		bitset := utils.CountryBitset{a, b, c, d}
		countries := bitset.ToCountries()
		val, _ := countries.Value()
		var countries2 utils.Countries
		countries2.Scan(val)
		assertCountriesEq(t, countries, countries2)
	})
}
