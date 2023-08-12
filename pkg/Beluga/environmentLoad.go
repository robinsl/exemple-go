package Beluga

import "github.com/joho/godotenv"

func LoadDotEnv() error {
	if err := godotenv.Load("configs/.env"); err != nil {
		return err
	}

	return nil
}
