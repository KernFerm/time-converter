package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var fallbackZones = []string{"UTC", "America/New_York", "Asia/Tokyo", "Europe/London"}

func loadTimezones(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("⚠️ Could not load timezones.json. Using fallback list.")
		return fallbackZones
	}
	var raw []string
	err = json.Unmarshal(data, &raw)
	if err != nil {
		fmt.Println("⚠️ Invalid JSON. Using fallback.")
		return fallbackZones
	}

	valid := []string{}
	seen := map[string]bool{}
	for _, tz := range raw {
		if _, err := time.LoadLocation(tz); err == nil && !seen[tz] {
			valid = append(valid, tz)
			seen[tz] = true
		}
	}
	return valid
}

func convertTime(input, from, to string) string {
	fromLoc, err1 := time.LoadLocation(from)
	toLoc, err2 := time.LoadLocation(to)
	if err1 != nil || err2 != nil {
		return "❌ Invalid timezones"
	}

	parsed, err := time.ParseInLocation("03:04 PM", strings.ToUpper(input), fromLoc)
	if err != nil {
		return "❌ Invalid time format. Use HH:MM AM/PM (e.g., 01:45 PM)"
	}

	return fmt.Sprintf("🕓 %s in %s is 🕒 %s", input, to, parsed.In(toLoc).Format("03:04 PM"))
}

func main() {
	_ = loadTimezones("timezones.json")
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🌐 Bubbles Timezone Converter")
	fmt.Println("👉 Format time as HH:MM AM/PM (e.g., 01:30 PM)")
	fmt.Println("💡 Type 'q' at any prompt to quit")
	fmt.Println("────────────────────────────────────────────")

mainloop:
	for {
		fmt.Print("\n⏰ Enter time (HH:MM AM/PM): ")
		inputTime, _ := reader.ReadString('\n')
		inputTime = strings.TrimSpace(inputTime)
		if strings.EqualFold(inputTime, "q") {
			break mainloop
		}

		fmt.Print("🌍 Enter source timezone (e.g., America/New_York): ")
		fromZone, _ := reader.ReadString('\n')
		fromZone = strings.TrimSpace(fromZone)
		if strings.EqualFold(fromZone, "q") {
			break mainloop
		}

		fmt.Print("🌎 Enter target timezone (e.g., Europe/London): ")
		toZone, _ := reader.ReadString('\n')
		toZone = strings.TrimSpace(toZone)
		if strings.EqualFold(toZone, "q") {
			break mainloop
		}

		result := convertTime(inputTime, fromZone, toZone)
		fmt.Println(result)
		fmt.Println("────────────────────────────────────────────")
	}

	fmt.Println("\n👋 Thanks for using Bubbles Timezone Converter! Goodbye!")
	fmt.Println("🫧 Press ENTER to close...")
	fmt.Scanln() // Keeps window open if double-clicked
}
