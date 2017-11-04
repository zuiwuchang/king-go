/*K3XsxSalt 是古典加密的 一個 加密 組件 用於 快速 將數據 編碼爲 不太容易 被解開的 密文
Salt 顯示了 隨機 撒鹽 以使相同數據 隨機 編碼爲不同的 數據 從而 避免一些 差量分析的逆向 行爲
3Xsx 顯示了 編碼 將 通過 三步 運算

1	產生 3個 隨機的 鹽 salts=[0,1,2] 並將 鹽值+一個固定值 Salt
2	對 待加密的每個 字節 進行 Xsx 三步運算
3	將 步驟1產生的鹽 salts 寫入輸出數據頭部 將加密後的數據 寫到 鹽後

Xsx
	//step 0
	b = b xor (salts[0] + Salt) //異或
	//step 1
	swapbit(b,0,7)	//交換 字節 第0bit 與 第7bit 的 值
	swapbit(b,1,6)
	swapbit(b,2,5)
	swapbit(b,3,4)
	b = b shl-loop (salts[1] + Salt) // shr-loop 左移 但 越位的數據 不捨棄 而是作爲最低位的 補位數據
	//step 2
	b = b xor (salts[2] + Salt)

*/
package classical

import (
	"crypto/rand"
	"errors"
	"github.com/zuiwuchang/king-go/crypto"
)

type K3XsxSalt struct {
	//鹽 基值
	Salt byte
}

//返回 鹽長度
func (k *K3XsxSalt) SaltLen() int {
	return 3
}

//返回 加密後 密文長度
func (k *K3XsxSalt) EncryptionLen(srcLen int) int {
	return srcLen + k.SaltLen()
}

//返回 解密後 密文長度
func (k *K3XsxSalt) DecryptionLen(srcLen int) int {
	return srcLen - k.SaltLen()
}

//加密
func (k *K3XsxSalt) EncryptionEx(src []byte, dist []byte) error {
	n := len(src)
	if k.EncryptionLen(len(src)) > len(dist) {
		return errors.New("dist buffer is too small,use EncryptionLen(len(src) get dist buffer len")
	}
	saltLen := k.SaltLen()
	_, err := rand.Read(dist[0:saltLen])
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		dist[saltLen+i] = k.encryption(src[i], dist[:saltLen])
	}
	return nil
}

//加密
func (k *K3XsxSalt) Encryption(b []byte) ([]byte, error) {
	n := len(b)
	saltLen := k.SaltLen()

	out := make([]byte, n+saltLen, n+saltLen)
	_, err := rand.Read(out[0:saltLen])
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		out[saltLen+i] = k.encryption(b[i], out[:saltLen])
	}
	return out, nil
}

//解密
func (k *K3XsxSalt) DecryptionEx(src []byte, dist []byte) error {
	n := len(src)
	if n < 3 {
		return errors.New("not encryption data")
	}
	if k.DecryptionLen(len(src)) > len(dist) {
		return errors.New("dist buffer is too small,use EncryptionLen(len(src) get dist buffer len")
	}

	saltLen := k.SaltLen()
	n -= saltLen
	for i := 0; i < n; i++ {
		dist[i] = k.decryption(src[i+saltLen], src[:saltLen])
	}
	return nil
}

//解密
func (k *K3XsxSalt) Decryption(b []byte) ([]byte, error) {
	n := len(b)
	if n < 3 {
		return nil, errors.New("not encryption data")
	}
	saltLen := k.SaltLen()
	n -= saltLen
	out := make([]byte, n, n)
	for i := 0; i < n; i++ {
		out[i] = k.decryption(b[i+saltLen], b[:saltLen])
	}
	return out, nil
}
func (k *K3XsxSalt) decryption(b byte, salts []byte) byte {
	bs := crypto.Byte{Data: b}
	//step 2
	bs.Xor(salts[2] + k.Salt)

	//step 1
	bs.ShrLoop(uint(salts[1] + k.Salt))
	bs.SwapBit(0, 7)
	bs.SwapBit(1, 6)
	bs.SwapBit(2, 5)
	bs.SwapBit(3, 4)

	//step 0
	bs.Xor(salts[0] + k.Salt)
	return bs.Data
}
func (k *K3XsxSalt) encryption(b byte, salts []byte) byte {
	bs := crypto.Byte{Data: b}
	//step 0
	bs.Xor(salts[0] + k.Salt)
	//step 1
	bs.SwapBit(0, 7)
	bs.SwapBit(1, 6)
	bs.SwapBit(2, 5)
	bs.SwapBit(3, 4)
	bs.ShlLoop(uint(salts[1] + k.Salt))

	//step 2
	bs.Xor(salts[2] + k.Salt)
	return bs.Data
}
