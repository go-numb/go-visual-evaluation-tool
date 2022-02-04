package modules

import "testing"

func TestCreateCSV(t *testing.T) {
	if err := CreateCSV("../data/images"); err != nil {
		t.Error(err)
	}
}
