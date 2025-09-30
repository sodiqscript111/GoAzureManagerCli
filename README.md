# GoAzureManagerCli

## What’s this?  
**GoAzureManagerCli** is a lightweight command-line tool for querying Azure Log Analytics in real time, built in Go.  
It lets you run KQL queries against your Azure workspace, export results to JSON, and expose Prometheus metrics — making it a handy tool for developers, DevOps engineers, and SREs working with Azure logs.  

---

## Features ✨  
- Run custom **KQL queries** against Azure Log Analytics 🔍  
- Query time ranges with simple flags (e.g., `5m`, `1h`, `24h`) ⏳  
- Export results to **JSON** 📄  
- Automatically track **404 errors** and increment Prometheus counters 📈  
- Expose metrics endpoint at `http://localhost:8080/metrics` for scraping ⚡  

---

## Getting Started 🚀  

### Prerequisites  
- Go 1.18+ installed  
- Azure credentials set as environment variables:  
  ```bash
  export AZURE_CLIENT_ID=your-client-id
  export AZURE_CLIENT_SECRET=your-client-secret
  export AZURE_TENANT_ID=your-tenant-id
  export AZURE_WORKSPACE_ID=your-workspace-id

Build and Run
go run main.go


The server will start and expose metrics at http://localhost:8080/metrics.

CLI Options ⚙️

-query → Custom KQL query (default:
"AppTraces | where TimeGenerated > ago(24h) | project TimeGenerated, Message")

-timespan → Time range for the query (e.g., 5m, 1h, 24h)

-output → Output format (default: json)

Example:

go run main.go -query "AppRequests | where ResultCode == 500" -timespan 1h -output json

API / Metrics 📊

Prometheus metrics are available at:

http://localhost:8080/metrics


Example metric output:

# HELP app_errors_total The total number of 404 errors observed
# TYPE app_errors_total counter
app_errors_total 12

Example Run 💻
# Run with defaults
go run main.go

# Output
Table: PrimaryResult
Column: TimeGenerated (datetime)
Column: Message (string)
Row: [2025-09-30T12:00:00Z "Request failed with 404"]
Row: [2025-09-30T12:05:00Z "Request completed successfully"]

Logs exported to logs.json
Metrics exposed at http://localhost:8080/metrics
