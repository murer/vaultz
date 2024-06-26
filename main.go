package main

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/murer/vaultz/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
)

const F_PUB = 0644
const F_PRIV = 0600

var Config = &packet.Config{
	DefaultHash:            crypto.SHA256,
	DefaultCipher:          packet.CipherAES256,
	DefaultCompressionAlgo: packet.CompressionZLIB,
	CompressionConfig: &packet.CompressionConfig{
		Level: packet.BestCompression,
	},
	RSABits: 512,
}

func Validate(c bool, msg string, v ...any) {
	if !c {
		log.Panicf(msg, v...)
	}
}

func GetBaseFile(filename string) string {
	base := os.Getenv("VAULTZ_BASE")
	if base == "" {
		log.Panic("VAULTZ_BASE is required")
	}
	ret := filepath.Join(base, filename)
	dir := filepath.Dir(ret)
	os.MkdirAll(dir, os.ModePerm)
	return ret
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func SHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

func GetBlob(filename string) string {
	hash := SHA256([]byte(filename))
	return GetBaseFile(filepath.Join("gen/blob", fmt.Sprintf("%s.secret.txt", hash)))
}

func ArmorIn(writer io.Writer, blockType string) io.WriteCloser {
	ret, err := armor.Encode(writer, blockType, nil)
	Check(err)
	return ret
}

func GenerateKeyPair(name string) {
	kp, err := openpgp.NewEntity(name, name, fmt.Sprintf("%s@any", name), Config)
	Check(err)
	log.Printf("Generating key %s: %s\n", name, kp.PrimaryKey.KeyIdString())
	file := GetBaseFile(fmt.Sprintf("pubkey/%s.pubkey.txt", name))
	pub, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, F_PUB)
	Check(err)
	defer pub.Close()
	(func() {
		apub := ArmorIn(pub, openpgp.PublicKeyType)
		defer apub.Close()
		kp.Serialize(apub)
	})()
	pub.Write([]byte{10})

	file = GetBaseFile("gen/privkey/privkey.txt")
	priv, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, F_PRIV)
	Check(err)
	defer priv.Close()
	(func() {
		apriv := ArmorIn(priv, openpgp.PrivateKeyType)
		defer apriv.Close()
		kp.SerializePrivate(apriv, Config)
	})()
	priv.Write([]byte{10})
}

func ReadKey(filename string) *openpgp.Entity {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	Check(err)
	defer file.Close()
	lst, err := openpgp.ReadArmoredKeyRing(file)
	util.Check(err)
	for _, v := range lst[0].Identities {
		v.SelfSignature.PreferredSymmetric = []uint8{uint8(packet.CipherAES256)}
		id, ok := s2k.HashToHashId(crypto.SHA256)
		util.Assert(!ok, "hash not found")
		v.SelfSignature.PreferredHash = []uint8{id}
	}
	return lst[0]
}

func ReadPubKeys() openpgp.EntityList {
	dir := GetBaseFile("pubkey")
	files, err := os.ReadDir(dir)
	Check(err)
	var ret openpgp.EntityList
	for _, file := range files {
		ret = append(ret, ReadKey(filepath.Join(dir, file.Name())))
	}
	return ret
}

func EncryptFile(filename string) {
	destfilename := GetBlob(filename)
	log.Printf("Encrypt %s: %s", filename, destfilename)
	pubkeys := ReadPubKeys()
	privkey := ReadKey(GetBaseFile("gen/privkey/privkey.txt"))
	file, err := os.OpenFile(filename, os.O_RDONLY, F_PRIV)
	Check(err)
	defer file.Close()
	destfile, err := os.OpenFile(destfilename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, F_PUB)
	Check(err)
	defer destfile.Close()
	(func() {
		adestfile := ArmorIn(destfile, "PGP MESSAGE")
		defer adestfile.Close()
		writer, err := openpgp.Encrypt(adestfile, pubkeys, privkey, nil, Config)
		Check(err)
		defer writer.Close()
		io.Copy(writer, file)
	})()
	destfile.Write([]byte{10})
}

func DecryptFile(filename string) {
	srcfilename := GetBlob(filename)
	log.Printf("Decrypt %s: %s", filename, srcfilename)
	pubkeys := ReadPubKeys()
	var kr openpgp.EntityList
	kr = append(kr, ReadKey(GetBaseFile("gen/privkey/privkey.txt")))
	srcfile, err := os.OpenFile(srcfilename, os.O_RDONLY, F_PUB)
	Check(err)
	defer srcfile.Close()
	areader, err := armor.Decode(srcfile)
	Check(err)
	if areader.Type != "PGP MESSAGE" {
		log.Panicf("Wrong block: %s\n", areader.Type)
	}
	msg, err := openpgp.ReadMessage(areader.Body, kr, nil, Config)
	Check(err)
	Validate(msg.IsEncrypted, "Message it not encrypted")
	Validate(!msg.IsSymmetricallyEncrypted, "Message is symmetrically encrypted")
	Validate(msg.IsSigned, "Message is not signed")
	Validate(len(pubkeys.KeysById(msg.SignedBy.PublicKey.KeyId)) > 0, "Message is not signed by any known key: %s", msg.SignedBy.PublicKey.KeyIdString())

	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, F_PRIV)
	Check(err)
	defer file.Close()
	io.Copy(file, msg.LiteralData.Body)
}

// ****************************************

type Command interface {
	GetName() string
	PrepareFlags(flags *flag.FlagSet)
	GetFlagSet() *flag.FlagSet
	Run()
}

type BaseCommand struct {
	Name    string
	Cmds    map[string]Command
	FlagSet *flag.FlagSet
}

func (me *BaseCommand) GetName() string {
	return me.Name
}

func (me *BaseCommand) PrepareFlags(flags *flag.FlagSet) {
	me.FlagSet = flags
}

func (me *BaseCommand) Run() {
	log.Println("aaaa")
}

func (me *BaseCommand) GetFlagSet() *flag.FlagSet {
	return me.FlagSet
}

// ****************************************

type HelpCommand struct {
	BaseCommand
}

func (me *HelpCommand) Run() {
	me.FlagSet.Output().Write([]byte{10})
	for _, cmd := range me.Cmds {
		cmd.GetFlagSet().Usage()
		os.Stdout.Write([]byte{10, 10})
	}
}

// ****************************************

type KeygenCommand struct {
	BaseCommand
	FlagName *string
}

func (me *KeygenCommand) Run() {
	GenerateKeyPair(*me.FlagName)
}

func (me *KeygenCommand) PrepareFlags(flags *flag.FlagSet) {
	me.FlagName = flags.String("name", "", "Key name")
	me.FlagSet = flags
}

// ****************************************

type EncryptCommand struct {
	BaseCommand
	FlagFile *string
}

func (me *EncryptCommand) Run() {
	EncryptFile(*me.FlagFile)
}

func (me *EncryptCommand) PrepareFlags(flags *flag.FlagSet) {
	me.FlagFile = flags.String("file", "", "File to be encrypted")
	me.FlagSet = flags
}

// ****************************************

type DecryptCommand struct {
	BaseCommand
	FlagFile *string
}

func (me *DecryptCommand) Run() {
	DecryptFile(*me.FlagFile)
}

func (me *DecryptCommand) PrepareFlags(flags *flag.FlagSet) {
	me.FlagFile = flags.String("file", "", "File to be decrypted")
	me.FlagSet = flags
}

// ****************************************

func createCommands() map[string]Command {
	ret := make(map[string]Command)
	(func(cmds []Command) {
		for _, cmd := range cmds {
			cmd.PrepareFlags(flag.NewFlagSet(cmd.GetName(), flag.ExitOnError))
			ret[cmd.GetName()] = cmd
		}
	})([]Command{
		&HelpCommand{BaseCommand: BaseCommand{"help", ret, nil}},
		&KeygenCommand{BaseCommand: BaseCommand{"keygen", ret, nil}},
		&EncryptCommand{BaseCommand: BaseCommand{"encrypt", ret, nil}},
		&DecryptCommand{BaseCommand: BaseCommand{"decrypt", ret, nil}},
	})
	return ret
}

func handleCommands(args []string) {
	commands := createCommands()
	subcommand := "help"
	if len(args) >= 2 {
		subcommand = args[1]
		args = args[2:]
	}
	command := commands[subcommand]
	if command == nil {
		log.Panicf("Wrong command: %s, try to use help", subcommand)
	}
	command.GetFlagSet().Parse(args)
	command.Run()
}

func main() {
	handleCommands(os.Args)
}
