package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "パスワード (必須)")
	flag.Parse()

	if *password == "" {
		fmt.Println("使用方法:")
		fmt.Println("go run cmd/hash_password.go -password=yourpassword")
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("パスワードハッシュ化エラー:", err)
	}

	fmt.Printf("パスワード: %s\n", *password)
	fmt.Printf("ハッシュ: %s\n", string(hashedBytes))
}