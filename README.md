# GCP Info Tool

A simple command-line utility to display Google Cloud Platform project information.

## Overview

This tool fetches and prints the following GCP information for the currently configured `gcloud` environment:
- Google Cloud Project ID
- Google Cloud Project Number
- Google Cloud Region Name

## Prerequisites

Before using this tool, you must have the [Google Cloud SDK (`gcloud`)](https://cloud.google.com/sdk/docs/install) installed and configured on your system. This means you should be authenticated:
```bash
gcloud auth login
gcloud auth application-default login
```
And ideally, have a default project and region configured:
```bash
gcloud config set project YOUR_PROJECT_ID
gcloud config set compute/region YOUR_REGION
```
The tool relies on `gcloud` to source this information. If `gcloud` is not found or not configured, the tool will output "N/A" for the respective fields and print errors to stderr.

## Usage

### From Pre-compiled Release

1.  Download the appropriate binary for your system (Linux or macOS) from the GitHub Releases page for this repository.
2.  Make the binary executable: `chmod +x ./gcp_info-linux-amd64` (or the binary you downloaded).
3.  Run the tool: `./gcp_info-linux-amd64`

Expected output:
```
google_cloud_project: your-project-id
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
2.  Build the tool:
    ```bash
    go build ./gcp_info.go
    ```
    This will create an executable named `gcp_info` (or `gcp_info.exe` on Windows) in the current directory.
3.  Run the built tool:
    ```bash
    ./gcp_info
    ```

## Releases

Pre-compiled binaries for Linux (amd64) and macOS (amd64, arm64) are automatically generated and attached to GitHub Releases whenever a new version tag (e.g., `v1.0.0`, `v1.1.0-beta.1`) is pushed to the repository.

## Contributing

Contributions are welcome! Please feel free to open an issue to discuss a new feature or bug, or submit a pull request with your changes.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
