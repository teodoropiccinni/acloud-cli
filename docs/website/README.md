# Aruba Cloud CLI Documentation Website

This website is built using [Docusaurus](https://docusaurus.io/), a modern static website generator.

## Installation

```bash
cd docs/website
npm install
```

## Local Development

```bash
npm start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

## Build

```bash
npm run build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

## Versioning

To create a new version of the documentation:

```bash
npm run version <version>
```

For example:
```bash
npm run version 1.0.0
```

This will:
- Create a new version directory in `versioned_docs/version-<version>/`
- Create a new version entry in `versions.json`
- Copy the current docs to the new version directory

## Deployment

The website is automatically deployed to GitHub Pages when a release is created via the GitHub Actions workflow.

