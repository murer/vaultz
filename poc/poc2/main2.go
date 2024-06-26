package main

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/murer/vaultz/pgp"
	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/ssh"
)

func parsePrivateKey(sshPrivateKey []byte) (*rsa.PrivateKey, error) {
	privateKey, err := ssh.ParseRawPrivateKey(sshPrivateKey)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)

	if !ok {
		return nil, fmt.Errorf("Only RSA keys are supported right now, got: %s", reflect.TypeOf(privateKey))
	}

	return rsaKey, nil
}

func SSHPrivateKeyToPGP(sshPrivateKey []byte, name string, comment string, email string) (*openpgp.Entity, error) {
	key, err := parsePrivateKey(sshPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private ssh key: %w", err)
	}

	// Let's make keys reproducible
	timeNull := time.Unix(0, 0)

	gpgKey := &openpgp.Entity{
		PrimaryKey: packet.NewRSAPublicKey(timeNull, &key.PublicKey),
		PrivateKey: packet.NewRSAPrivateKey(timeNull, key),
		Identities: make(map[string]*openpgp.Identity),
	}
	uid := packet.NewUserId(name, comment, email)
	isPrimaryID := true
	gpgKey.Identities[uid.Id] = &openpgp.Identity{
		Name:   uid.Id,
		UserId: uid,
		SelfSignature: &packet.Signature{
			CreationTime:              timeNull,
			SigType:                   packet.SigTypePositiveCert,
			PubKeyAlgo:                packet.PubKeyAlgoRSA,
			Hash:                      crypto.SHA256,
			IsPrimaryId:               &isPrimaryID,
			FlagsValid:                true,
			FlagSign:                  true,
			FlagCertify:               true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &gpgKey.PrimaryKey.KeyId,
		},
	}
	err = gpgKey.Identities[uid.Id].SelfSignature.SignUserId(uid.Id, gpgKey.PrimaryKey, gpgKey.PrivateKey, nil)
	if err != nil {
		return nil, err
	}

	return gpgKey, nil
}

func main2() {
	dat, err := os.ReadFile("gen/ssh/key")
	util.Check(err)

	entity, err := SSHPrivateKeyToPGP(dat, "any", "any", "any")
	util.Check(err)
	fmt.Printf("entity: %v\n", entity)

	kp := pgp.KeyFromEntity(entity)
	fmt.Printf("kp: %#v\n", kp)

	os.WriteFile("gen/ssh/pgp.priv.txt", []byte(kp.ExportPrivArmored()), 0600)
	os.WriteFile("gen/ssh/pgp.pub.txt", []byte(kp.ExportPubArmored()), 0644)

	kr := pgp.KeyRingCreate(kp)

	crypt := pgp.EncryptString("plain", kp, kr)
	os.WriteFile("gen/ssh/encrypted.bin", crypt, 0644)

	plain := pgp.DecryptString(crypt, kr, kr)
	fmt.Printf("Plain: %s\n", plain)

	fmt.Println("Done")

}
