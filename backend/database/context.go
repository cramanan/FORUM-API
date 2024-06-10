package database

import (
	"context"
	"fmt"
)

func RegisterUser(ctx context.Context, nickname, age, gender, firstname, lastname, email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := "INSERT INTO Users (id, nickname, age, gender, firstname, lastname, email, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, query, GenerateUUID(), nickname, age, gender, firstname, lastname, email, hashedPassword)
	if err != nil {
		return nil
	}
	err = tx.Commit()
	if err != nil {
		return nil
	}

	fmt.Println("User registered successfully")

	return nil
}
