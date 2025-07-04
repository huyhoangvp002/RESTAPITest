package main

import (
	"RESTAPITest/api"
	db "RESTAPITest/db/sqlc"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	//  Lấy biến môi trường DB_SOURCE
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable not set")
	}

	//  Mở kết nối database
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("cannot open database:", err)
	}

	// Test kết nối DB
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	fmt.Println("Connected to the database!")

	//Tạo store (Querier)
	store := db.New(conn)

	//Tạo API server
	server := api.NewServer(store)

	//Chạy server trên cổng 8080
	addr := ":8080"
	fmt.Printf("Server is running at %s\n", addr)
	err = server.Start(addr)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
