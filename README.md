# GCP Info Tool

A simple command-line utility to display Google Cloud Platform project information using the Google Cloud APIs.

## Overview

This tool fetches and prints the following GCP information for a **specified project ID** passed as a command-line argument:
- Google Cloud Project ID (this will be the ID you provide)
- Google Cloud Project Number (fetched based on the provided Project ID using Resource Manager API)
- Google Cloud Region Name

Region information is determined first by checking a project-specific location label (`cloud.googleapis.com/location`) on the provided project; if not found or empty, it falls back to retrieving region information from the Compute Engine API.

## Prerequisites

Before using this tool, you must have authentication set up for Google Cloud. The tool uses [Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/application-default-credentials) to authenticate with Google Cloud APIs.

You can set up ADC by:
```bash
gcloud auth application-default login
```

Or by setting the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to a service account key file:
```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-key.json"
```

The authenticated user or service account must have the necessary IAM permissions (e.g., `resourcemanager.projects.get`) to access the project ID you provide.

## Usage

### From Pre-compiled Release

1.  Download the appropriate binary for your system (Linux or macOS) from the GitHub Releases page for this repository.
2.  Make the binary executable: `chmod +x ./gcp_info-linux-amd64` (or the binary you downloaded).
3.  Run the tool, providing the target project ID as an argument:
    ```bash
    ./gcp_info-linux-amd64 your-project-id-here
    ```

Expected output:
```
Fetching info for project: your-project-id-here
google_cloud_project: your-project-id-here
google_cloud_project_number: 123456789012
google_cloud_region_name: us-central1
```
If any value cannot be determined, it will show `N/A`.

### Building from Source

1.  Clone the repository:
    ```bash
    git clone <repository_clone_url> # Replace with actual repository clone URL
    cd <repository_directory_name> # Replace with actual repository directory name
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Build the tool:
    ```bash
    go build ./gcp_info.go
    ```
    This will create an executable named `gcp_info` (or `gcp_info.exe` on Windows) in the current directory.
4.  Run the built tool, providing the target project ID as an argument:
    ```bash
    ./gcp_info your-project-id-here
    ```

## Releases

Pre-compiled binaries for Linux (amd64) and macOS (amd64, arm64) are automatically generated and attached to GitHub Releases whenever a new version tag (e.g., `v1.0.0`, `v1.1.0-beta.1`) is pushed to the repository.

## Required APIs

The following Google Cloud APIs must be enabled in your project:
- Resource Manager API
- Compute Engine API

You can enable these APIs using the Google Cloud Console or with the following commands:
```bash
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable compute.googleapis.com
```

## Contributing

Contributions are welcome! Please feel free to open an issue to discuss a new feature or bug, or submit a pull request with your changes.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
