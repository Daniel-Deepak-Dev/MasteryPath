package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DashboardSubSkill struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Mastery float64 `json:"mastery"`
}

type DashboardSkill struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Category       string              `json:"category"`
	OverallMastery float64             `json:"overall_mastery"`
	SubSkills      []DashboardSubSkill `json:"sub_skills"`
}

func main() {
	baseURL := "http://localhost:8080/api"
	fmt.Println("=== 🧪 Starting MasteryPath API Verification ===")

	// Test Dashboard
	pass := testDashboard(baseURL + "/skills/dashboard")

	if pass {
		fmt.Println("\n✅ ALL TESTS PASSED")
	} else {
		fmt.Println("\n❌ TESTS FAILED")
	}
}

func testDashboard(endpoint string) bool {
	fmt.Printf("\nTesting Endpoint: %s\n", endpoint)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		log.Printf("❌ Connection Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Unexpected Status Code: %d\n", resp.StatusCode)
		return false
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var skills []DashboardSkill
	if err := json.Unmarshal(body, &skills); err != nil {
		log.Printf("❌ JSON Parse Error: %v\n", err)
		return false
	}

	if len(skills) == 0 {
		log.Println("❌ No skills returned (Expected at least Salesforce Development)")
		return false
	}

	// Validate Salesforce Skill
	foundSF := false
	for _, s := range skills {
		if s.Name == "Salesforce Development" {
			foundSF = true
			fmt.Printf("  ✓ Found Master Skill: %s (Mastery: %.1f%%)\n", s.Name, s.OverallMastery)

			// Validate Sub-skills exist
			expectedSubs := map[string]bool{"Apex": false, "Triggers": false, "LWC": false}

			// Mark found subs
			for _, sub := range s.SubSkills {
				if _, exists := expectedSubs[sub.Name]; exists {
					expectedSubs[sub.Name] = true
					fmt.Printf("    ✓ Found Sub-skill: %s (Mastery: %.1f%%)\n", sub.Name, sub.Mastery)
				}
			}

			// Check if all expected subs found
			allSubsFound := true
			for name, found := range expectedSubs {
				if !found {
					fmt.Printf("    ❌ Missing Sub-skill: %s\n", name)
					allSubsFound = false
				}
			}
			if !allSubsFound {
				return false
			}
		}
	}

	if !foundSF {
		log.Println("❌ 'Salesforce Development' skill not found in response")
		return false
	}

	return true
}
