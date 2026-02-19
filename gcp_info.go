package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./gcp_info <project_id>")
		os.Exit(1)
	}
	targetProjectID := os.Args[1]

	// Initialize display variables
	projectIDStr := targetProjectID
	projectNumberStr := "N/A"
	regionNameStr := "N/A"

	fmt.Fprintf(os.Stdout, "Fetching info for project: %s\n", targetProjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get project information using Resource Manager API
	projectInfo, err := getProjectInfo(ctx, targetProjectID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting project info for %s: %v\n", projectIDStr, err)
	} else {
		projectNumberStr = fmt.Sprintf("%d", projectInfo.ProjectNumber)

		// Check for project-specific location label
		if locationLabel, ok := projectInfo.Labels["cloud.googleapis.com/location"]; ok && locationLabel != "" {
			regionNameStr = locationLabel
			fmt.Fprintf(os.Stderr, "Info: Used project-specific location label for region: %s\n", regionNameStr)
		} else {
			fmt.Fprintf(os.Stderr, "Info: Project-specific location label for %s is not set. Trying default region.\n", projectIDStr)
			
			// Try to get default region from Compute Engine API
			defaultRegion, err := getDefaultRegion(ctx, targetProjectID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting default region: %v\n", err)
			} else if defaultRegion == "" {
				fmt.Fprintln(os.Stderr, "Warning: Default region not found. Region will be N/A.")
			} else {
				regionNameStr = defaultRegion
				fmt.Fprintf(os.Stderr, "Info: Used default region: %s\n", regionNameStr)
			}
		}
	}

	// Print final formatted output
	fmt.Printf("google_cloud_project: %s\n", projectIDStr)
	fmt.Printf("google_cloud_project_number: %s\n", projectNumberStr)
	fmt.Printf("google_cloud_region_name: %s\n", regionNameStr)
}

// getProjectInfo retrieves project information using the Resource Manager API
func getProjectInfo(ctx context.Context, projectID string) (*cloudresourcemanager.Project, error) {
	rmService, err := cloudresourcemanager.NewService(ctx, option.WithScopes(cloudresourcemanager.CloudPlatformScope))
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource Manager service: %w", err)
	}

	projectsService := rmService.Projects
	project, err := projectsService.Get(projectID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get project info: %w", err)
	}

	return project, nil
}

// getDefaultRegion attempts to determine the default region for the project
func getDefaultRegion(ctx context.Context, projectID string) (string, error) {
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeReadonlyScope))
	if err != nil {
		return "", fmt.Errorf("failed to create Compute service: %w", err)
	}

	// First try to find any compute instance and use its region
	instancesList, err := computeService.Instances.AggregatedList(projectID).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to list instances: %w", err)
	}

	// Look for any instance and extract region from its zone
	for _, instancesScopedList := range instancesList.Items {
		if instancesScopedList.Instances != nil && len(instancesScopedList.Instances) > 0 {
			// Get the first instance found
			instance := instancesScopedList.Instances[0]
			// Zone is in format: "projects/PROJECT_ID/zones/ZONE_NAME"
			// or just "ZONE_NAME", extract the zone name
			zoneParts := strings.Split(instance.Zone, "/")
			zoneName := zoneParts[len(zoneParts)-1]
			
			// Zone names are typically in the format "us-central1-a"
			// Extract region by removing the last part (e.g., "-a")
			parts := strings.Split(zoneName, "-")
			if len(parts) >= 2 {
				// Remove the last part and rejoin
				region := strings.Join(parts[:len(parts)-1], "-")
				fmt.Fprintf(os.Stderr, "Info: Used region from existing compute instance: %s\n", region)
				return region, nil
			}
		}
	}

	// Fallback: try to get the project's metadata
	project, err := computeService.Projects.Get(projectID).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get project metadata: %w", err)
	}

	// Look for default region in metadata
	if project.CommonInstanceMetadata != nil && project.CommonInstanceMetadata.Items != nil {
		for _, item := range project.CommonInstanceMetadata.Items {
			if item.Key == "google-compute-default-region" && item.Value != nil {
				fmt.Fprintf(os.Stderr, "Info: Used region from project metadata: %s\n", *item.Value)
				return *item.Value, nil
			}
		}
	}

	// Final fallback: try to get zone information and extract region from it
	zonesList, err := computeService.Zones.List(projectID).Context(ctx).Do()
	if err == nil && len(zonesList.Items) > 0 {
		// Zone names are typically in the format "us-central1-a"
		// Extract region by removing the last part (e.g., "-a")
		zoneName := zonesList.Items[0].Name
		parts := strings.Split(zoneName, "-")
		if len(parts) >= 2 {
			// Remove the last part and rejoin
			region := strings.Join(parts[:len(parts)-1], "-")
			fmt.Fprintf(os.Stderr, "Info: Used region from available zone: %s\n", region)
			return region, nil
		}
	}

	return "", nil
}