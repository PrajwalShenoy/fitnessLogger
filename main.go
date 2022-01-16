package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	gosseract "github.com/otiai10/gosseract/v2"
)

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func read_json_config(path string) map[string]string {
	if path == "" {
		path = "ftp_config.json"
	}
	jsonFile, err := os.Open(path)
	check_err(err)
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	var config map[string]string
	json.Unmarshal([]byte(byteValue), &config)
	return config
}

func make_ftp_connection(ip, port, username, password string) *ftp.ServerConn {
	c, err := ftp.Dial(ip + ":" + port)
	check_err(err)
	err = c.Login(username, password)
	check_err(err)
	return c
}

func create_directory_structure(path string) {
	year, month, day := time.Now().Date()
	parent_dir := path + "/" + fmt.Sprintf("%d%02d%02d", year, month, day)
	err := os.MkdirAll(parent_dir, os.ModePerm)
	check_err(err)
	sc_dir := parent_dir + "/Screenshots"
	pic_dir := parent_dir + "/Pictures"
	err = os.MkdirAll(sc_dir, os.ModePerm)
	check_err(err)
	err = os.MkdirAll(pic_dir, os.ModePerm)
	check_err(err)
}

func pull_screenshots(conn *ftp.ServerConn, source_path, dest_path string) {
	listDir, err := conn.NameList(source_path)
	check_err(err)
	year, month, day := time.Now().Date()
	screenshot_path := dest_path + "/" + fmt.Sprintf("%d%02d%02d/Screenshots", year, month, day)
	for _, screenshot := range listDir {
		if strings.Contains(screenshot, fmt.Sprintf("Screenshot_%d%02d%02d", year, int(month), day)) {
			response, err := conn.Retr(source_path + "/" + screenshot)
			check_err(err)
			buffer, err := ioutil.ReadAll(response)
			check_err(err)
			response.Close()
			os.WriteFile(screenshot_path+"/"+screenshot, buffer, 0644)
		}
	}
}

func pull_photos(conn *ftp.ServerConn, source_path, dest_path string) {
	listDir, err := conn.NameList(source_path)
	check_err(err)
	year, month, day := time.Now().Date()
	photo_path := dest_path + "/" + fmt.Sprintf("%d%02d%02d/Pictures", year, month, day)
	for _, photo := range listDir {
		if strings.Contains(photo, fmt.Sprintf("IMG_%d%02d%02d", year, month, day)) {
			response, err := conn.Retr(source_path + "/" + photo)
			check_err(err)
			buffer, err := ioutil.ReadAll(response)
			check_err(err)
			response.Close()
			os.WriteFile(photo_path+"/"+photo, buffer, 0644)
		}
	}
}

func process_screenshots(path string) (string, string, string, string) {
	var fiber, carbs, fats, protein string
	year, month, day := time.Now().Date()
	screenshot_path := path + "/" + fmt.Sprintf("%d%02d%02d/Screenshots", year, month, day)
	listDir, err := ioutil.ReadDir(screenshot_path)
	check_err(err)
	for _, file := range listDir {
		client := gosseract.NewClient()
		client.SetImage(screenshot_path + "/" + file.Name())
		text, _ := client.Text()
		if strings.Contains(text, "fibre Consumed") {
			fiber = extract_value(text, "Fiber")
		} else if strings.Contains(text, "carbs Consumed") {
			carbs = extract_value(text, "Carbs")
		} else if strings.Contains(text, "fats Consumed") {
			fats = extract_value(text, "Fats")
		} else if strings.Contains(text, "protein Consumed") {
			protein = extract_value(text, "Protein")
		} else {
			fmt.Println("Screenshot does not fit in the preset templates")
		}
	}
	fmt.Printf("Fiber %s; Carbs %s; Fats %s; Protein %s\n", fiber, carbs, fats, protein)
	return fiber, carbs, fats, protein
}

func write_to_csv(path_to_csv, fiber, carbs, fats, protein string) {
	file, err := os.OpenFile(path_to_csv, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check_err(err)
	defer file.Close()
	year, month, date := time.Now().Date()
	string_date := fmt.Sprintf("%d/%02d/%02d", year, month, date)
	text_for_csv := fmt.Sprintf("%s,%s,%s,%s,%s\n", string_date, fiber, carbs, fats, protein)
	_, err = file.WriteString(text_for_csv)
	check_err(err)
}

func extract_value(text, macro string) string {
	var to_ret string
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "/") {
			to_ret = strings.Split(line, "/")[0]
		}
	}
	return to_ret
}

func main() {
	config := read_json_config("ftp_config.json")
	connection := make_ftp_connection(config["ftp_ip"],
		config["ftp_port"],
		config["username"],
		config["password"])
	create_directory_structure(config["local_store_path"])
	pull_screenshots(connection, config["screenshot_path"], config["local_store_path"])
	pull_photos(connection, config["photos_path"], config["local_store_path"])
	fiber, carbs, fats, protein := process_screenshots(config["local_store_path"])
	write_to_csv(config["csv_path"], fiber, carbs, fats, protein)
	connection.Quit()
}
