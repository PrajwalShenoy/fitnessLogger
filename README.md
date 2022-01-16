# Fitness Logger

## Installing pre-requisites
Install tesseract using 

```bash
sudo apt install tesseract-ocr
sudo apt install libtesseract-dev
```

Install gosseract using
```bash
go get  github.com/otiai10/gosseract/v2
```

Install ftp using
```bash
go get github.com/jlaffaye/ftp
```

## Pre-requisites for application to run
Create a file called `ftp_config.json
Create the JSON structure as seen below
```json
{
    "ftp_ip": IP_ADDR,
    "ftp_port": PORT,
    "username": FTP_USERNAME,
    "password": FTP_PASSWORD,
    "screenshot_path": PATH_TO_SCREENSHOT_ON_PHONE,
    "photos_path": PATH_TO_DCMI_PHOTOS_ON_PHONE,
    "local_store_path": DIR_PATH_TO_STORE_PHOTOS_AND_SCREENSHOTS,
    "csv_path": PATH_TO_CSV_FILE
}
```
Create a CSV file with the following as the first line
```
Date,Fiber,Carbohydrates,Fats,Protein
<new line>
```

## Running the application for testing purposes
```bash
go run main.go
```

## Creating a binary for execution
```
go build -o fitness_tracker
```