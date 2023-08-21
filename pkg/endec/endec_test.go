package endec

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var key []byte

func init() {
	key, _ = base64.StdEncoding.DecodeString("NK6Ucg8hCXN9oPZFojqcAqYEY2DETZzt6oKrGmZwdOU=")
}

func TestAesCBC(t *testing.T) {
	raw := []byte("fuck")
	enc, _ := AesCBCEncrypt(raw, key, key[0:16])
	fmt.Printf("enc: %v\n", enc)
	dec, _ := AesCBCDecrypt(enc, key, key[0:16])
	fmt.Printf("dec: %v\n", dec)
}

func TestAesCFB(t *testing.T) {
	raw := []byte("fuck")
	enc, _ := AesCFBEncrypt(raw, key, key[0:16])
	fmt.Printf("enc: %v\n", enc)
	dec, _ := AesCFBDecrypt(enc, key, key[0:16])
	fmt.Printf("dec: %v\n", dec)
}

func TestHk4eAbilityHashCode(t *testing.T) {
	hashCode := Hk4eAbilityHashCode("Avatar_Ayato_ExtraAttack")
	fmt.Printf("Avatar_Ayato_ExtraAttack hashCode: %v\n", hashCode)
	hashCode = Hk4eAbilityHashCode("Avatar_Ayato_ExtraAttack_CreateBullet")
	fmt.Printf("Avatar_Ayato_ExtraAttack_CreateBullet hashCode: %v\n", hashCode)
}
