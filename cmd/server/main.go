package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"example.com/mod/internal/db"
	genproto "example.com/mod/internal/genproto"
	grpcimpl "example.com/mod/internal/grpcimpl"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database connection
	dsn := "user:your_password@tcp(localhost:3306)/dbname"
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Ensure the database is reachable
	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize the database wrapper
	database := db.NewDatabase(conn)

	// Set up listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// path to the SQL file in the migrations folder
	createDecisionTable := "migrations/001_create_decisions_table.sql"
	indexOnDecisionTable := "migrations/002_index_decisions_table.sql"
	insertDecisionTable := "migrations/003_insert_decisions_table.sql"

	// Read the SQL file for creating the table
	createDecTable, err := os.ReadFile(createDecisionTable)
	if err != nil {
		log.Fatalf("Error reading create decision migration script: %v", err)
	}

	// Execute the SQL to create the table
	_, err = conn.Exec(string(createDecTable))
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// If successful, print a success message
	fmt.Println("Table created successfully")

	// Read the SQL file for inserting data
	indexDecTable, err := os.ReadFile(indexOnDecisionTable)
	if err != nil {
		log.Fatalf("Error reading index decision migration script: %v", err)
	}

	fmt.Println("string(insertDecTable)", string(indexDecTable))
	// Execute the SQL to insert data
	_, err = conn.Exec(string(indexDecTable))
	if err != nil {
		log.Fatalf("Error creating index on table: %v", err)
	}

	// If successful, print a success message
	fmt.Println("Index created successfully")

	// Read the SQL file for inserting data
	insertDecTable, err := os.ReadFile(insertDecisionTable)
	if err != nil {
		log.Fatalf("Error reading insert decision migration script: %v", err)
	}

	// Execute the SQL to insert data
	_, err = conn.Exec(string(insertDecTable))
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	// If successful, print a success message
	fmt.Println("Data inserted successfully")

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the ExploreService with the server using the generated Register function
	genproto.RegisterExploreServiceServer(s, &grpcimpl.ExploreServer{DB: database})

	// Start server
	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
