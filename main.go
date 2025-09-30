package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var errorCount = promauto.NewCounter(prometheus.CounterOpts{Name: "app_errors_total"})

func main() {
	query := flag.String("query", "AppTraces | where TimeGenerated > ago(24h) | project TimeGenerated, Message", "KQL query")
	output := flag.String("output", "json", "Output format: json")
	timespan := flag.String("timespan", "24h", "Time range (e.g., 5m, 1h, 24h)")
	flag.Parse()

	// Parse timespan
	duration, err := time.ParseDuration(*timespan)
	if err != nil {
		fmt.Printf("Invalid timespan: %v\n", err)
		os.Exit(1)
	}

	os.Setenv("AZURE_CLIENT_ID", "c31360fd-9470-46dd-863b-ee68a4fe6e20")
	os.Setenv("AZURE_CLIENT_SECRET", "l_K8Q~UUF2JDFwl_nITPlfglKsGs0.6pbmwwocQ~")
	os.Setenv("AZURE_TENANT_ID", "7ba79ee2-d48e-4c72-bccb-e31872af6ea0")
	workspaceID := "c252bad2-f34b-40c3-bfba-f17f1d0425f0"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Printf("Credential error: %v\n", err)
		os.Exit(1)
	}

	logsClient, err := azquery.NewLogsClient(cred, nil)
	if err != nil {
		fmt.Printf("Client error: %v\n", err)
		os.Exit(1)
	}

	res, err := logsClient.QueryWorkspace(
		context.TODO(),
		workspaceID,
		azquery.Body{
			Query:    to.Ptr(*query),
			Timespan: to.Ptr(azquery.NewTimeInterval(time.Now().Add(-duration), time.Now())),
		},
		nil,
	)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		os.Exit(1)
	}

	// Parse and increment Prometheus metric (count errors)
	if len(res.Tables) > 0 {
		errorCountVal := 0.0
		for _, table := range res.Tables {
			for _, row := range table.Rows {
				if msg, ok := row[1].(string); ok && strings.Contains(msg, "404") {
					errorCountVal++
				}
			}
		}
		errorCount.Add(errorCountVal)
	}

	// CLI output
	if len(res.Tables) > 0 {
		for _, table := range res.Tables {
			fmt.Println("Table:", *table.Name)
			for _, col := range table.Columns {
				fmt.Printf("Column: %s (%s)\n", *col.Name, *col.Type)
			}
			for _, row := range table.Rows {
				fmt.Printf("Row: %v\n", row)
			}
		}
	} else {
		fmt.Println("No results returned")
	}

	// JSON export
	if *output == "json" {
		jsonData, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Printf("JSON error: %v\n", err)
			os.Exit(1)
		}
		if err := ioutil.WriteFile("logs.json", jsonData, 0644); err != nil {
			fmt.Printf("File write error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Logs exported to logs.json")
	}

	// Prometheus HTTP endpoint
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)
	fmt.Println("Metrics exposed at http://localhost:8080/metrics")
	select {} // Keep running
}
