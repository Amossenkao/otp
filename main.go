package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"otp/utils"

	"github.com/spf13/cobra"
)

//go:embed utils/*
var util embed.FS

// Global variables

var password string
var services []Service
var customEnv = Env{}

func init() {
	jsonFile, err := util.Open("utils/env.json")
	
	if err != nil {
		log.Fatalf("Error reading utils/services.json: %v", err)
	}
	defer jsonFile.Close()

	byteResult, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteResult, &customEnv)
	services = customEnv.Services
	password = customEnv.Password

}

type Env struct {
	Services []Service `json:"services"`
	Password string    `json:"password"`
}

type Service struct {
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Secret    string `json:"secret"`
}


func getSecret(service string) (string, string, error) {
	service = strings.ToLower(service)
	for i := 0; i < len(services); i++ {

		svc := services[i]
		if svc.Name == service || svc.ShortName == service {
			return svc.Secret, svc.Name, nil
		}

	}
	return "", "", fmt.Errorf("unrecognized service: %s", service)

}

func main() {

	var list bool
	var rootCmd = &cobra.Command{Use: "TOTP Generator", Example: "otp get [service]", Run: func(cmd *cobra.Command, args []string) {
		if list {
			if !utils.Authenticate(password) {
				log.Fatal("Authentication failed")
			} else {

				count := 1
				for i := 0; i < len(services); i++ {
					svc := services[i]
					fmt.Printf("%d. Name: %s <==> Short Name: %s\n", count, svc.Name, svc.ShortName)
					count++
				}
			}
		}
	}}

	var getCmd = &cobra.Command{
		Use:     "get [service]",
		Short:   "Generate OTP for a service",
		Long:    "The get command takes the name of a spervice and retrives the TOTP code for that service",
		Example: "otp get github, otp get ig",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			service := args[0]
			secret, _, err := getSecret(service)
			if err != nil {
				log.Fatal(err)
			}
			for i := 1; i <= 3; i++ {
				totpCode, remaining, err := utils.GetOtp(secret)
				if err == nil {
					fmt.Printf("TOTP generated and copied to the clipboard: %s\n", totpCode)
					fmt.Printf("Valid for the next %d seconds\n", remaining)
					 if i < 3 {
						time.Sleep(time.Duration(remaining) * time.Second)
					 }
				}
			}
		},
	}

	var secretCmd = &cobra.Command{
		Use:     "secret [service]",
		Short:   "Get the secret for a service",
		Long:    "The secret command returns the secret for a sprcified service after successful authentication",
		Example: "otp secret gmail, otp get secret fb",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			service := args[0]
			if !utils.Authenticate(password) {
				log.Fatal("Authentication failed")
			}

			secret, serviceName, err := getSecret(service)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Secret for %s: %s\n", serviceName, secret)
		},
	}

	// Add subcommands
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(secretCmd)
	rootCmd.Flags().BoolVarP(&list, "list", "l", false, "List all services")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
