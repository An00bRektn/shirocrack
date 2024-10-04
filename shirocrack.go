package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"strconv"
)

var salt []byte
var iterationCount int
var b64Hash string

// https://github.com/apache/shiro/blob/eef749eb994ab6d5f975f9bd96481e4d7d427ac0/crypto/hash/src/main/java/org/apache/shiro/crypto/hash/SimpleHash.java#L386
func hashWithSaltAndIterations(bytes []byte, salt []byte, iterations int) []byte {
	digest := sha512.New()

	// Apply salt before the first hash
	digest.Reset()
	if len(salt) != 0 {
		digest.Write(salt)
	}
	digest.Write(bytes)
	hashed := digest.Sum(nil)

	// Perform remaining iterations
	for i := 1; i < iterations; i++ {
		digest.Reset()
		digest.Write(hashed)
		hashed = digest.Sum(nil)
	}

	return hashed
}

// Function to compare the hash of a plaintext password with the expected base64-encoded hash
func hash(plaintextPassword string) bool {
	hashed := hashWithSaltAndIterations([]byte(plaintextPassword), salt, iterationCount)

	// Encode the final hash to base64 for comparison
	hashBase64 := base64.StdEncoding.EncodeToString(hashed)

	if hashBase64 == b64Hash {
		fmt.Printf("[+] Found match: %s:%s\n", hashBase64, plaintextPassword)
		return true
	}
	return false
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("[i] usage: ./shirocrack HASH wordlist.txt")
		return
	}

	// this is the hash
	mcfString := os.Args[1]

	// Split the MCF string to extract iteration count, salt, and the expected hash
	// $shiro1$SHA-512$ITERATIONS$B64_SALT$B64_HASH
	mcf := strings.Split(mcfString, "$")
	iterationCount, _ = strconv.Atoi(mcf[3]) // make sure they're not malformed
	b64Salt := mcf[4]
	b64Hash = mcf[5]

	// Decode the base64-encoded salt
	salt, _ = base64.StdEncoding.DecodeString(b64Salt)

	fileBytes, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	cracked := false

	// Process the dictionary file line by line
	lines := strings.Split(string(fileBytes), "\n")
	for _, line := range lines {
		if hash(line) {
			cracked = true
			break
		}
	}

	if cracked {
		fmt.Println("[+] Success!")
	} else {
		fmt.Println("[!] Failed to crack hash. :(")
	}
	
}
