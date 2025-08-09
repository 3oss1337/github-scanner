Github Scanner
used for scanning repos for .env files and any API keys or AWS keys
to use clone the repo and in cmd/server, go run main.go
then using powershell type 
    Invoke-RestMethod -Method Post -Uri http://localhost:8080/scan -ContentType "application/json" -Body '{"owner":"yourAccountName","repo":"yourRepoName"}' | ConvertTo-Json -Depth 5
and it will return null if you don't have any vulnerabilities in your repo, and if there it will return the key/env
you can also use curl
    curl -X POST http://localhost:8080/scan -H "Content-Type: application/json" -d '{"owner":"yourAccountName","repo":"yourRepoName"}'
