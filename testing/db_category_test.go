package testing

import (
	"log"
	"testing"

	. "slugquest.com/backend/crud"
)

func TestGetCategory(t *testing.T) bool {
	cat, bol, erro := GetCatId(50)
	if !bol {
		log.Println("TestGetCat(): Get Cat ID(): cat id not found")
	}
	if erro != nil {
		log.Printf("TestGetCat(): Get Cat ID() #1: %v", erro)
		return false
	}

	if cat.CatID != 50 {
		log.Println("TestGcat(): found wrong cat")
		return false
	}

	cat, bol, erro = GetCatId(-5)
	if bol {
		log.Printf("TestGetCat(): Get Cat ID():  find catad")
		return false
	}
	if erro != nil {
		log.Printf("TestGetCat(): Get Cat ID() #2: %v", erro)
		return false
	}

	return true
}
