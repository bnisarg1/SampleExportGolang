package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type User struct {
	SubjectID string
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", "localhost", 5502, "username", "password", "hsdp_pg")

	dbConn, errDbConn := sql.Open("postgres", psqlInfo)
	if errDbConn != nil {
		panic(errDbConn)
	}

	defer dbConn.Close()

	// Ensure the database is reachable
	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	// Define the buffered channel
	bufferSize := 10000
	ch := make(chan User, bufferSize)

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Note: uncomment the following code block to read data from the database
	/*// Define the start and end times
	startTime := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 5, 30, 0, 0, 0, 0, time.UTC)

	// Calculate the duration for each time slot
	duration := endTime.Sub(startTime) / 4

	// Start the read goroutine
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(slot int) {
			defer wg.Done()
			slotStartTime := startTime.Add(time.Duration(slot) * duration)
			slotEndTime := slotStartTime.Add(duration)
			readFromDB(dbConn, ch, slotStartTime, slotEndTime)
		}(i)
	}*/

	// Read data from the CSV file and send to the channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		readFromCSVFile(dbConn, ch)
	}()

	/*// Start a goroutine to close the channel after all readers are done
	go func() {
		wg.Wait()
		close(ch)
	}()
	var counter int
	for user := range ch {
		fmt.Println("User: ", user.SubjectID)
		counter++
	}
	fmt.Println("Total number of users: ", counter)
	*/

	// Process the data from the channel
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range ch {
				processAndInsertUser(dbConn, user)
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("Finished processing all users.")

}

func readFromCSVFile(conn *sql.DB, ch chan User) {

	log.Printf("Reading data from csv file")
	// Open the CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	// Specify the target column index (0-based)
	targetColumnIndex := 1 // Change this to the desired column index

	// Skip header row (if present)
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		fmt.Println("Error reading header:", err)
		return
	}

	// Read all records and store them in the buffer
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // Reached end of file
		}
		if err != nil {
			fmt.Println("Error reading record:", err)
			return
		}

		// Extract the target column value
		if strings.EqualFold(record[4], "") {
			ch <- User{SubjectID: record[targetColumnIndex]}
		}
	}

}

func processAndInsertUser(db *sql.DB, user User) {
	fmt.Printf("Processing user: %v\n", user)
	externalid := uuid.New()
	query := "INSERT INTO mappingTableName (userid, externalid, type) VALUES ('%s', '%s' , 'IAM')"
	stmt := fmt.Sprintf(query, user.SubjectID, externalid)
	_, err := db.Exec(stmt)
	if err != nil {
		log.Printf("Failed to insert user %d: %v", user.SubjectID, err)
		return
	}

	fmt.Printf("Inserted user %v with subject_id %v into another_table\n", user.SubjectID, externalid)

}

// Read data from the database for a specific time range and send to the channel
func readFromDB(db *sql.DB, ch chan<- User, startTime, endTime time.Time) {
	log.Printf("Reading data from %v to %v", startTime, endTime)
	query := "SELECT subjectid FROM  someTableName WHERE lastmodified >= '%s' AND lastmodified < '%s'"
	stmt := fmt.Sprintf(query, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	rows, err := db.Query(stmt)
	if err != nil {
		log.Printf("Failed to query from %v to %v: %v", startTime, endTime, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.SubjectID); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		ch <- user
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error in rows: %v", err)
	}
}
