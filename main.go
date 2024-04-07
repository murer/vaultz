package main

import (
	"bytes"
	"crypto"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

const F_PUB = 0644
const F_PRIV = 0600

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

func Config() *packet.Config {
	return &packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: packet.BestCompression,
		},
		RSABits: 512,
	}
}

func ArmorIn(writer io.Writer, blockType string) io.WriteCloser {
	ret, err := armor.Encode(writer, blockType, nil)
	Check(err)
	return ret
}

func ArmorInBytes(data []byte) string {
	buf := new(bytes.Buffer)
	func() {
		writer := ArmorIn(buf, "PGP MESSAGE")
		defer writer.Close()
		writer.Write(data)
	}()
	return buf.String()
}

func GenerateKeyPair(name string) {
	kp, err := openpgp.NewEntity(name, name, fmt.Sprintf("%s@any", name), Config())
	Check(err)
	log.Printf("Generating key %s: %s\n", name, kp.PrimaryKey.KeyIdString())
	file := GetBaseFile(fmt.Sprintf("pubkey/%s.pubkey.txt", name))
	writer, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, F_PUB)
	Check(err)
	(func() {
		defer writer.Close()
		(func() {
			awriter := ArmorIn(writer, openpgp.PublicKeyType)
			defer awriter.Close()
			kp.PrimaryKey.Serialize(awriter)
		})()
		writer.Write([]byte{10})
	})()

	file = GetBaseFile("gen/privkey/privkey.txt")
	writer, err = os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, F_PRIV)
	Check(err)
	(func() {
		defer writer.Close()
		(func() {
			awriter := ArmorIn(writer, openpgp.PrivateKeyType)
			defer awriter.Close()
			kp.PrivateKey.Serialize(awriter)
		})()
		writer.Write([]byte{10})
	})()
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
	log.Printf("Command: %s, args: %s\n", subcommand, args)
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
