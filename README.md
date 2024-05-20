# SampleExportGolang
Demo app for Exporting Data from CSV file/Some Table to specified table

#Goal
For Golang Learning purpose created the following sample app which make use of Goroutine and buffered channel 
- Following topics covered
  1) Buffered channel
  2) Executing multiple Goroutine to read and write
  3) DB connection setup
  4) waitsync group
     

#Steps to run the following 
- fork or clone the repo
- go mod tidy
- execute via intellj IDE or one can execute it by running go run main.go
  

#Code walk through
- Commented code are used to switch between CSV or Select Query from DB
- we are reading the second column subject id(From data.csv file) if alternativeuseruuid  is empty
- Then we are passing it to the channel where multiple gorutine will read the subjectid and insert data to the table
- To create multiple goroutine for loop is used

- For Reading From DB
- Uncommnent the code and comment the CSV reading flow
- Data read based on the start time and end time  -  since requirement was to run mulitple goroutie to read  the data from db, so start time and end time is divided into slots
  and passed to the goroutine
- processAndInsertUser: will read the subjectid from the channel and insert it
