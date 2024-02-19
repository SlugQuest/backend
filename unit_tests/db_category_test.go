package unit_tests

import (
	"testing"

	. "slugquest.com/backend/crud"
)

func TestGetCategory(t *testing.T) {
	cat, found, erro := GetCatId(50)
	if erro != nil {
		t.Errorf("TestGetCat(): Get Cat ID() #1: %v", erro)
	}

	if !found {
		t.Error("TestGetCat(): Get Cat ID(): cat id not found")
	}

	if cat.CatID != 50 {
		t.Error("TestGetCat(): found wrong cat")
	}

	cat, found, erro = GetCatId(-5)
	if erro != nil {
		t.Errorf("TestGetCat(): Get Cat ID() #2: %v", erro)
	}

	if found {
		t.Error("TestGetCat(): bad find cat, found -5")
	}
}
