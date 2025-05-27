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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./gcp_info <project_id>")
		os.Exit(1)
	}
	targetProjectID := os.Args[1]

	// Initialize display variables
	projectIDStr := targetProjectID // Display the user-provided ID
	projectNumberStr := "N/A"
	regionNameStr := "N/A" // Region is typically not project-specific in the same way, 
	                        // but we'll fetch it in the context of the project if possible later.
	                        // For now, it's N/A as per this step's focus.

	fmt.Fprintf(os.Stdout, "Fetching info for project: %s\n", targetProjectID)

	// Get Project Number for the targetProjectID
	// targetProjectID is projectIDStr
	if projectIDStr != "" { // Should always be true due to earlier os.Args check, but good for clarity
		projectNumberFetched, err := executeGcloudCommand("projects", "describe", projectIDStr, "--format=value(projectNumber)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting project number for %s: %v\n", projectIDStr, err)
		} else if projectNumberFetched == "" {
			fmt.Fprintf(os.Stderr, "Warning: Project number from gcloud for %s was empty; displaying as N/A.\n", projectIDStr)
		} else {
			projectNumberStr = projectNumberFetched
		}
	}

	// Get Region Information
	// Attempt 1: Project-specific location label
	// Use projectIDStr which is targetProjectID
	regionLabelFetched, errLabel := executeGcloudCommand("projects", "describe", projectIDStr, "--format=value(labels.cloud.googleapis.com/location)")
	if errLabel == nil && regionLabelFetched != "" {
		regionNameStr = regionLabelFetched
		fmt.Fprintf(os.Stderr, "Info: Used project-specific location label for region: %s\n", regionNameStr)
	} else {
		if errLabel != nil {
			// This error might occur if the project doesn't exist or if the label truly causes a command failure.
			// gcloud often returns empty string for non-existent label with exit code 0, but if --format errors, it could be non-zero.
			fmt.Fprintf(os.Stderr, "Info: Could not get project-specific location label for %s (may not be set or project query failed): %v. Trying local default region.\n", projectIDStr, errLabel)
		} else if regionLabelFetched == "" {
			fmt.Fprintf(os.Stderr, "Info: Project-specific location label for %s is empty. Trying local default region.\n", projectIDStr)
		}

		// Attempt 2: Local gcloud default region (fallback)
		// This is only attempted if the project-specific label wasn't found or was empty.
		localRegionFetched, errLocalRegion := executeGcloudCommand("config", "get-value", "compute/region")
		if errLocalRegion != nil {
			fmt.Fprintf(os.Stderr, "Error getting local default region: %v\n", errLocalRegion)
		} else if localRegionFetched == "" {
			// This warning is important if no region could be determined at all.
			fmt.Fprintln(os.Stderr, "Warning: Local default region from gcloud was empty. Region will be N/A.")
		} else {
			// Only use local default if project-specific one wasn't successfully set.
			if regionNameStr == "N/A" { 
				regionNameStr = localRegionFetched
				fmt.Fprintf(os.Stderr, "Info: Used local default gcloud region: %s\n", regionNameStr)
			}
		}
	}

	// Print final formatted output
	fmt.Printf("google_cloud_project: %s\n", projectIDStr)
	fmt.Printf("google_cloud_project_number: %s\n", projectNumberStr)
	fmt.Printf("google_cloud_region_name: %s\n", regionNameStr)
}
