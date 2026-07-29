package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	var fname, lname string
	var number int
	fmt.Print("First name:")
	fmt.Scan(&fname)
	fmt.Print("Last name:")
	fmt.Scan(&lname)
	fmt.Print("Number:")
	fmt.Scan(&number)

	fmt.Printf("F: %s\n", fname)
	fmt.Printf("L: %s\n", lname)
	fmt.Printf("C: %d\n", number)

	conn, err := pgx.Connect(context.Background(), "postgres://admin:password@localhost:5432/userdata")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var user string
	conn.QueryRow(context.Background(),
		"SELECT current_user").Scan(&user)
	fmt.Println(user)

	_, err = conn.Exec(context.Background(), `INSERT INTO public.user_info (fname, lname, number) VALUES ($1, $2, $3)`, fname, lname, number)
	if err != nil {
		log.Fatal(err)
	}

}
