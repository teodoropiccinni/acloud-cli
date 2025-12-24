# Documentation Deployment Guide

This guide explains how the documentation is deployed to GitHub Pages.

## Automatic Deployment

The documentation is automatically deployed when:

1. **A new release is published** - The GitHub Actions workflow (`docs-release.yml`) triggers on release events
2. **Manual trigger** - You can manually trigger the workflow from the Actions tab

## How It Works

1. When a release is created (or manually triggered), the workflow:
   - Extracts the version from the release tag (removes 'v' prefix if present)
   - Creates a versioned copy of the docs using `npm run version <version>`
   - Builds the Docusaurus site
   - Deploys to GitHub Pages

2. The versioned documentation will be available at:
   - Latest: `https://arubacloud.github.io/acloud-cli/`
   - Versioned: `https://arubacloud.github.io/acloud-cli/docs/<version>/`

## Manual Deployment

If you need to deploy manually:

```bash
cd docs/website

# Install dependencies (first time only)
npm install

# Create a new version (if needed)
npm run version 1.0.0

# Build the site
npm run build

# The build output is in the build/ directory
# You can test it locally with:
npm run serve
```

## GitHub Pages Setup

1. Go to your repository settings
2. Navigate to "Pages" section
3. Set source to "GitHub Actions"
4. The workflow will automatically deploy to GitHub Pages

## Versioning

- Each release creates a new version of the documentation
- Users can switch between versions using the version dropdown
- The "Next" version shows the latest development docs

## Troubleshooting

- If deployment fails, check the GitHub Actions logs
- Ensure GitHub Pages is enabled in repository settings
- Verify the `baseUrl` in `docusaurus.config.js` matches your repository name

