package testing

import (
	"testing"

	. "slugquest.com/backend/crud"
)

func TestGetCategory(t *testing.T) {
	cat, bol, erro := GetCatId(50)
	if !bol {
		t.Error("TestGetCat(): Get Cat ID(): cat id not found")
	}
	if erro != nil {
		t.Errorf("TestGetCat(): Get Cat ID() #1: %v", erro)
	}

	if cat.CatID != 50 {
		t.Error("TestGcat(): found wrong cat")
	}

	cat, bol, erro = GetCatId(-5)
	if bol {
		t.Error("TestGetCat(): Get Cat ID():  find catad")
	}
	if erro != nil {
		t.Errorf("TestGetCat(): Get Cat ID() #2: %v", erro)
	}
}
