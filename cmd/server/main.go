package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/3oss1337/github-scanner/githubapi"
	"github.com/3oss1337/github-scanner/models"
	"github.com/3oss1337/github-scanner/scanner"
	"github.com/3oss1337/github-scanner/utils"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		var req models.ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Owner) == "" || strings.TrimSpace(req.Repo) == "" {
			http.Error(w, "owner and repo are required", http.StatusBadRequest)
			return
		}

		token := utils.GetGithubToken()
		client := githubapi.NewClient(token)

		results := make([]models.ScanResponse, 0, 128)

		var walk func(path string)
		walk = func(path string) {
			entries := client.ListRepoContents(req.Owner, req.Repo, path)
			if entries == nil {
				return
			}
			for _, entry := range entries {
				if entry == nil || entry.Path == nil || entry.Type == nil {
					continue
				}
				switch *entry.Type {
				case "file":
					// Skip common binary or large non-text types
					lowerName := strings.ToLower(entry.GetName())
					if strings.HasSuffix(lowerName, ".png") || strings.HasSuffix(lowerName, ".jpg") ||
						strings.HasSuffix(lowerName, ".jpeg") || strings.HasSuffix(lowerName, ".gif") ||
						strings.HasSuffix(lowerName, ".pdf") || strings.HasSuffix(lowerName, ".zip") ||
						strings.HasSuffix(lowerName, ".tar") || strings.HasSuffix(lowerName, ".gz") {
						continue
					}
					file, err := client.GetFile(req.Owner, req.Repo, *entry.Path)
					if err != nil || file == nil {
						continue
					}
					content, _ := file.GetContent()
					scanResults := scanner.ScanFile(*entry.Path, content)
					for _, r := range scanResults {
						results = append(results, models.ScanResponse{File: *entry.Path, Type: r.Type, Match: r.Match})
					}
					if strings.HasSuffix(lowerName, ".env") || strings.Contains(lowerName, ".env") {
						if strings.TrimSpace(content) != "" {
							results = append(results, models.ScanResponse{File: *entry.Path, Type: ".env-like file", Match: "potential secrets present"})
						}
					}
				case "dir":
					walk(*entry.Path)
				default:
				}
			}
		}

		walk("")

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
