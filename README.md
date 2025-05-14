# Image Syncer

A tool to sync container images from any public container registry to GitHub Container Registry (ghcr.io).

## Features

- Sync container images from any public registry to GitHub Container Registry (ghcr.io)
- Run as a GitHub Action with manual triggering
- Configurable source image and target organization
- Simple integration with GitHub Actions workflows

## How It Works

The Image Syncer:

1. Pulls the source container image from a public registry
2. Tags it for GitHub Container Registry
3. Authenticates with GitHub Container Registry
4. Pushes the image to GitHub Container Registry

## Why Use Image Syncer?

- **Simplified Dependency Management**: Keep copies of external container images in your own GitHub Container Registry
- **Improved Reliability**: Reduce dependency on external registries that might have downtime
- **Version Control**: Maintain specific versions of container images that work with your applications
- **Access Control**: Manage access to container images through GitHub's permissions system

## Usage as a GitHub Action

### Option 1: Using the Action Directly in Your Workflow

You can use this action directly in your workflow by referencing it:

```yaml
name: Sync Container Image

on:
  workflow_dispatch:
    inputs:
      source_image:
        description: 'Source container image to sync'
        required: true
        default: 'nginx:latest'

jobs:
  sync:
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: read
    
    steps:
      - name: Sync container image
        uses: ODearEvanHansen/image-syncer@main
        with:
          source_image: ${{ github.event.inputs.source_image }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
```

### Option 2: Setting Up a Complete Workflow

1. Add the GitHub Action workflow to your repository by creating a file at `.github/workflows/image-sync.yml`
2. Ensure your repository has the necessary permissions to write packages

### Manually Triggering the Action

1. Go to your repository on GitHub
2. Click on the "Actions" tab
3. Select the "Image Sync" workflow from the list
4. Click on "Run workflow"
5. Enter the required parameters:
   - **Source Image**: The container image to sync (e.g., `nginx:latest`, `ubuntu:20.04`)
   - **Target Organization** (optional): The target organization in GHCR (defaults to the repository owner)
6. Click "Run workflow"

### Example Workflow File

```yaml
name: Image Sync

on:
  workflow_dispatch:
    inputs:
      source_image:
        description: 'Source container image to sync (e.g., nginx:latest, ubuntu:20.04)'
        required: true
      target_org:
        description: 'Target organization in GHCR (default: current repository owner)'
        required: false
        default: ${{ github.repository_owner }}

jobs:
  sync:
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: read
    
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.19'

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build image-syncer
        run: |
          go build -o image-syncer ./cmd/image-syncer

      - name: Sync image
        env:
          GITHUB_ACTOR: ${{ github.actor }}
        run: |
          TARGET_ORG="${{ github.event.inputs.target_org }}"
          if [ -z "$TARGET_ORG" ]; then
            TARGET_ORG="${{ github.repository_owner }}"
          fi
          
          ./image-syncer \
            -source "${{ github.event.inputs.source_image }}" \
            -target-org "$TARGET_ORG" \
            -token "${{ secrets.GITHUB_TOKEN }}"
```

## Running Locally

### Prerequisites

- Go 1.19 or later
- Docker

### Building

```bash
go build -o image-syncer ./cmd/image-syncer
```

### Running

```bash
./image-syncer \
  -source "nginx:latest" \
  -target-org "your-github-username" \
  -token "your-github-token"
```

## License

This project is licensed under the terms of the license included in the repository.

## Testing

### Running Tests

```bash
go test -v ./pkg/syncer
```

### Testing the GitHub Action

You can test the GitHub Action by creating a workflow that uses the action and manually triggering it with a public container image.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request