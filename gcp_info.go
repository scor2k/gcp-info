package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func executeGcloudCommand(args ...string) (string, error) {
	cmd := exec.Command("gcloud", args...)
	output, err := cmd.Output() // cmd.Output() captures stdout
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Command executed but returned non-zero exit code
			return "", fmt.Errorf("gcloud command failed with: %s, stderr: %s", exitError, string(exitError.Stderr))
		} else if err == exec.ErrNotFound {
			// gcloud command not found
			return "", fmt.Errorf("gcloud command not found. Please ensure it is installed and in your PATH.")
		}
		// Other errors (e.g., I/O issues)
		return "", fmt.Errorf("failed to run gcloud command: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	projectIDStr := "N/A"
	projectNumberStr := "N/A"
	regionNameStr := "N/A"

    // ---- REFINED ASSIGNMENT LOGIC ----
    // Get Project ID
	projectIDFetched, err := executeGcloudCommand("config", "get-value", "project")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting project ID: %v\n", err)
	} else if projectIDFetched == "" {
		fmt.Fprintln(os.Stderr, "Warning: Project ID from gcloud was empty; displaying as N/A.")
	} else {
		projectIDStr = projectIDFetched
	}

	// Get Project Number
	if projectIDStr != "N/A" { // Only try if we have a valid project ID
		projectNumberFetched, err := executeGcloudCommand("projects", "describe", projectIDStr, "--format=value(projectNumber)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting project number for %s: %v\n", projectIDStr, err)
		} else if projectNumberFetched == "" {
			fmt.Fprintf(os.Stderr, "Warning: Project number from gcloud for %s was empty; displaying as N/A.\n", projectIDStr)
		} else {
			projectNumberStr = projectNumberFetched
		}
	}

	// Get Region
	regionNameFetched, err := executeGcloudCommand("config", "get-value", "compute/region")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting region: %v\n", err)
	} else if regionNameFetched == "" {
		fmt.Fprintln(os.Stderr, "Warning: Region from gcloud was empty; displaying as N/A.")
	} else {
		regionNameStr = regionNameFetched
	}

	// Print final formatted output
	fmt.Printf("google_cloud_project: %s\n", projectIDStr)
	fmt.Printf("google_cloud_project_number: %s\n", projectNumberStr)
	fmt.Printf("google_cloud_region_name: %s\n", regionNameStr)
}
