package main

import (
	"fmt"
)

func main() {
    encodedHash := "$argon2id$v=19$m=65536,t=3,p=2$1234567890abcdef$0987654321fedcba"
	var memory uint32
	var iterations uint32
	var parallelism uint8
	var saltHex string
	n, err := fmt.Sscanf(encodedHash, "$argon2id$v=19$m=%d,t=%d,p=%d$%s", &memory, &iterations, &parallelism, &saltHex)
	fmt.Printf("n: %d, err: %v\n", n, err)
	fmt.Printf("memory: %d, iterations: %d, parallelism: %d\n", memory, iterations, parallelism)
	fmt.Printf("saltHex: '%s'\n", saltHex)
}
