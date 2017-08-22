package classical

import (
	//"fmt"
	"testing"
)

func TestK3XsxSaltEx(t *testing.T) {
	k3 := K3XsxSalt{Salt: 1}

	str := "this is 測試 text"
	b := []byte(str)
	//fmt.Println("old : ", b)

	n := k3.EncryptionLen(len(b))
	enc := make([]byte, n)
	err := k3.EncryptionEx(b, enc)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println("enc : ", enc)

	n = k3.DecryptionLen(len(enc))
	dec := make([]byte, n)
	err = k3.DecryptionEx(enc, dec)
	if err != nil {
		t.Fatal(err)
	}

	if string(dec) != str {
		t.Fatal("Encryption not equal Decryption")
	}

	{
		str := ""
		b := []byte(str)
		//fmt.Println("old : ", b)

		n := k3.EncryptionLen(len(b))
		enc := make([]byte, n)
		err := k3.EncryptionEx(b, enc)
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println("enc : ", enc)

		n = k3.DecryptionLen(len(enc))
		dec := make([]byte, n)
		err = k3.DecryptionEx(enc, dec)
		if err != nil {
			t.Fatal(err)
		}

		if string(dec) != str {
			t.Fatal("Encryption not equal Decryption")
		}
	}
}
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
	{
		str := ""
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
}
