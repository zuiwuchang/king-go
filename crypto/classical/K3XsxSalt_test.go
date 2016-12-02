package classical

import (
	//"fmt"
	"testing"
)

func TestK3XsxSalt(t *testing.T) {
	k3 := K3XsxSalt{Salt: 1}

	str := "this is 測試 text"
	b := []byte(str)
	//fmt.Println("old : ", b)
	b, err := k3.Encryption(b)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println("enc : ", b)

	b, err = k3.Decryption(b)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != str {
		t.Fatal("Encryption not equal Decryption")
	}

}
