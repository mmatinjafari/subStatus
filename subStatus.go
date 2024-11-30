package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"io/ioutil"
	"encoding/json"
)

// ساختار برای ذخیره اطلاعات خروجی httpx
type HttpxResult struct {
	Host        string `json:"host"`
	StatusCode  int    `json:"status_code"`
}

// اجرای subfinder برای پیدا کردن زیر دامنه‌ها
func runSubfinder(domain string) {
	cmd := exec.Command("subfinder", "-d", domain, "-o", "subdomains.txt")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running subfinder: %v\n", err)
		return
	}
	fmt.Println("subfinder execution completed, subdomains saved to subdomains.txt.")
}

// اجرای httpx برای بررسی وضعیت HTTP زیر دامنه‌ها
func runHttpx() {
	// اجرای httpx و ذخیره خروجی به صورت JSON در output.json
	cmd := exec.Command("httpx", "-mc", "200,301,302", "-silent", "-json", "-o", "output.json", "-l", "subdomains.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running httpx: %v\nOutput: %s", err, string(output))
		return
	}
	fmt.Println("httpx execution completed, output saved to output.json.")
}

// پردازش نتایج خروجی httpx و ذخیره در valid_domains.txt
func processHttpxResults() {
	// خواندن فایل خروجی JSON
	data, err := ioutil.ReadFile("output.json")
	if err != nil {
		log.Printf("Error reading output.json: %v\n", err)
		return
	}

	// تجزیه داده‌های JSON
	var results []HttpxResult
	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Printf("Error parsing JSON data: %v\n", err)
		return
	}

	// باز کردن فایل valid_domains.txt برای نوشتن
	file, err := os.Create("valid_domains.txt")
	if err != nil {
		log.Printf("Error creating valid_domains.txt: %v\n", err)
		return
	}
	defer file.Close()

	// نوشتن نتایج در فایل valid_domains.txt
	for _, result := range results {
		// نوشتن هر زیر دامنه همراه با کد وضعیت HTTP
		line := fmt.Sprintf("%s, %d\n", result.Host, result.StatusCode)
		_, err := file.WriteString(line)
		if err != nil {
			log.Printf("Error writing to valid_domains.txt: %v\n", err)
			return
		}
	}
	fmt.Println("Valid domains saved to valid_domains.txt.")
}

func main() {
	// باز کردن فایل لاگ برای نوشتن خطاها
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v\n", err)
	}
	defer logFile.Close()

	// تنظیم لاگ برای نوشتن در فایل
	log.SetOutput(logFile)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run subStatus.go <domain>")
		return
	}

	domain := os.Args[1]

	// اجرای subfinder برای پیدا کردن زیر دامنه‌ها
	runSubfinder(domain)

	// اجرای httpx برای بررسی وضعیت HTTP با کدهای 200، 301 و 302
	runHttpx()

	// پردازش نتایج خروجی httpx و ذخیره در فایل valid_domains.txt
	processHttpxResults()
}
