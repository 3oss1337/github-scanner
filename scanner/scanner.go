package scanner

import (
	"regexp"
	"sync"

	"github.com/3oss1337/github-scanner/models"
)

var patterns = []models.Pattern{
	{Name: "AWS Access Key", Regex: `AKIA[0-9A-Z]{16}`},
	{Name: "AWS Secret Key", Regex: `(?i)aws_secret_access_key.*?[A-Za-z0-9/+=]{40}`},
	{Name: "Private Key", Regex: `-----BEGIN( RSA)? PRIVATE KEY-----`},
	{Name: "Generic Secret", Regex: `(?i)(SECRET|TOKEN|PASSWORD|API_KEY|DB_URI)[\s:=]+['\"]?.+['\"]?`},
}

func ScanFile(path string, content string) []models.ScanResult {
	var results []models.ScanResult
	var wg sync.WaitGroup
	resultsCh := make(chan models.ScanResult)
	for _, pattern := range patterns {
		wg.Add(1)
		go func(pattern models.Pattern) {
			defer wg.Done()
			re, err := regexp.Compile(pattern.Regex)
			if err != nil {
				return
			}
			matches := re.FindAllString(content, -1)
			for _, match := range matches {
				resultsCh <- models.ScanResult{
					Type:  pattern.Name,
					Match: match,
				}
			}
		}(pattern)
	}
	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	for result := range resultsCh {
		results = append(results, result)
	}
	return results
}
