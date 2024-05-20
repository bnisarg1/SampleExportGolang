# SampleExportGolang

A demo application for exporting data from a CSV file or a database table to a specified table, created for learning purposes in Golang. This project demonstrates the use of Goroutines and buffered channels.

## Features
The following topics are covered in this project:
1. Buffered channels
2. Executing multiple Goroutines for reading and writing
3. Database connection setup
4. Using the `sync.WaitGroup`

## Getting Started

### Prerequisites
- Go installed on your system (version 1.21.3 or later recommended)
- Access to a database 

### Installation
1. Fork or clone the repository:
   ```sh
   git clone https://github.com/yourusername/SampleExportGolang.git
   cd SampleExportGolang
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Execute the application:
   - Using IntelliJ IDEA or any Go-supported IDE
   - Or run it from the command line:
     ```sh
     go run main.go
     ```

## Code Walkthrough
- The code is designed to switch between reading from a CSV file or executing a SELECT query from a database.
- It reads the second column (`subjectid`) from `data.csv` if `alternativeuseruuid` is empty.
- The `subjectid` is then passed to a channel, where multiple Goroutines read the `subjectid` and insert data into the specified table.
- A loop is used to create multiple Goroutines.

### Reading from a CSV File
- By default, the code reads from a CSV file. If you want to read from a database, comment out the CSV reading code and uncomment the database reading code.

### Reading from a Database
- The data is read based on the `start_time` and `end_time`. 
- Since the requirement was to run multiple Goroutines to read data from the database, the `start_time` and `end_time` are divided into slots and passed to the Goroutines.
- The `processAndInsertUser` function reads the `subjectid` from the channel and inserts it into the table.

## Usage
1. Ensure your CSV file (`data.csv`) is in the correct format and located in the project directory.
2. If reading from a database, configure the database connection settings in the code.
3. Adjust the number of Goroutines and buffer sizes as needed in the code.

## Contributing
If you have suggestions for improvement or want to contribute, please fork the repository and create a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

This README file provides an overview of the project, instructions for setup and execution, and a brief code walkthrough. Feel free to improve it further based on your specific needs and updates to the project.
