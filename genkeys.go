package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ed25519"
)

func genKeys(conf Conf, configFile string) {
	randRead, randReader := rand.Read, io.Reader(nil)

	// psk is an optionally deterministically generated random byte slice
	psk := make([]byte, 32)
	if _, err := randRead(psk); err != nil {
		log.Fatal(err)
	}
	pskHex := hex.EncodeToString(psk)

	encryptSk := make([]byte, 32)
	if _, err := randRead(encryptSk); err != nil {
		log.Fatal(err)
	}
	encryptSkHex := hex.EncodeToString(encryptSk)

	signPk, signSk, err := ed25519.GenerateKey(randReader)
	if err != nil {
		log.Fatal(err)
	}
	signPkHex := hex.EncodeToString(signPk)
	signSkHex := hex.EncodeToString(signSk[0:32])

	fmt.Printf("\n\n--- Create a file named %s with only the lines relevant to your configuration ---\n\n\n", configFile)
	fmt.Printf("# Configuration for a client\n\n")
	fmt.Printf("Connect   = %q\t# Edit appropriately\n", conf.Connect)
	fmt.Printf("Psk       = %q\n", pskHex)
	fmt.Printf("SignPk    = %q\n", signPkHex)
	fmt.Printf("SignSk    = %q\n", signSkHex)
	fmt.Printf("EncryptSk = %q\n", encryptSkHex)

	fmt.Printf("\n\n")

	fmt.Printf("# Configuration for a server\n\n")
	fmt.Printf("Listen = %q\t# Edit appropriately\n", conf.Listen)
	fmt.Printf("Psk    = %q\n", pskHex)
	fmt.Printf("SignPk = %q\n", signPkHex)

	fmt.Printf("\n\n")

	fmt.Printf("# Hybrid configuration\n\n")
	fmt.Printf("Connect   = %q\t# Edit appropriately\n", conf.Connect)
	fmt.Printf("Listen    = %q\t# Edit appropriately\n", conf.Listen)
	fmt.Printf("Psk       = %q\n", pskHex)
	fmt.Printf("SignPk    = %q\n", signPkHex)
	fmt.Printf("SignSk    = %q\n", signSkHex)
	fmt.Printf("EncryptSk = %q\n", encryptSkHex)
}

func getPassword(prompt string) string {
	os.Stdout.Write([]byte(prompt))
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(password)
}
