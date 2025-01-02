package utils

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
)

// GetOtp generates a TOTP code and copies it to the clipboard

func GetOtp(secret string) (string, int, error) {
	
	currentTime := time.Now()
	totpCode, err := totp.GenerateCode(secret, currentTime)
	if err != nil {
		return "", 0, fmt.Errorf("error generating TOTP: %v", err)
	}

	err = clipboard.WriteAll(totpCode)
	if err != nil {
		return "", 0, fmt.Errorf("failed to copy TOTP to clipboard: %v", err)
	}
	return totpCode, getRemainingSeconds(currentTime), nil
}

// getRemainingSeconds returns the number of seconds remaining before the next TOTP code is generated

func getRemainingSeconds(currentTime time.Time) int {
	interval := 30
	remaining := interval - int(currentTime.Unix() % int64(interval))
	return remaining
}