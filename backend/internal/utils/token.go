package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

// GenerateDeviceCode generates a unique device identifier
func GenerateDeviceCode() string {
	return uuid.New().String()
}

// GenerateHex generates a random hex string of specified length
func GenerateHex(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateWToken generates the wToken with random hex prefix
func GenerateWToken() (string, error) {
	hexStr, err := GenerateHex(104)
	if err != nil {
		return "", err
	}

	// Build wToken with template
	wToken := fmt.Sprintf("0004_9D60C%sBaMraleAcDLVQnHi8CUzaKFfDnmVzEbM/ppBoLmcNUoD+RMN0TtCCCS+OWFpuRizcjcQdVTu1h5dlw7LZXICUs6rwa/bqdvPOV+bgaszYlLmil5s8hmCtA/mrxNrJzRgF8hFehAnWKlzTlErW1k83tAX3p5MotZVhW7gGfJhokkKpij1zqCv/F+kgiEySEWJSE8KUN+o7gcd86pDFiavjf+EUXBqNWzpEqgA+H1T8KfeYcC2YSEjtNQ14zQ0ers+MXcTgsugjQhyj7Qm1MbICgmZVUc+24+CtEfAhZ3VzGWMbYhtxzV0d59OVG82+gCURle2XNZVIdPGHLIdWa3/OxWZt8TgqxsFm+AbVfD8s5Cd/AOHR1uuSa4MavRVahMjFD2oyxWyy8xg2xPYk95BM4P1zZV57jdJ/kXEMBKJ+kT/NqVfct5rv+S9sLOAOsOuKsD7kSqWTUgb3Fq14wqEha5wZzBV4zyXXI6PQtA7ZYNirAytOmNeSt3vEgtbsmYiaEq1iztfWV3K7fmolWuUmMtXaB9EpjArTn7rnVIa9IzqoqJamG8fq/szkGU3rC0tzeidcxm7kgqqltRJpPNv4wmwbeVyFMWnNSdhrDwTl9zVowMNXgPUJSN1SbQfAGDgeYmkNxavGuZ5wrEmXM6AS4ZeXd/R/L01h7ffqpLSDw84dmbR7NkvXGXUYFAlsASx_fHw=", hexStr)

	return wToken, nil
}

// GenerateUCDE returns the observed login UCDE value.
func GenerateUCDE() (string, error) {
	return "t698", nil
}
