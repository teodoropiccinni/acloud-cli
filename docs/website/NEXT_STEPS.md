# Next Steps for Docusaurus Setup

## 1. Install Dependencies

First, install the Node.js dependencies:

```bash
cd docs/website
npm install
```

## 2. Test Locally

Start the development server to preview the documentation:

```bash
npm start
```

This will start a local server (usually at http://localhost:3000) where you can preview the documentation.

## 3. Customize the Site

### Update Logo and Favicon

Replace the placeholder files:
- `static/img/logo.svg` - Main logo (should be a square SVG)
- `static/img/favicon.ico` - Browser favicon

### Update Homepage Features

Edit `src/components/HomepageFeatures/index.js` to customize the three feature cards on the homepage.

### Customize Colors

Edit `src/css/custom.css` to change the color scheme.

## 4. Enable GitHub Pages

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Pages**
3. Under **Source**, select **GitHub Actions**
4. Save the settings

## 5. Test the Deployment Workflow

You can test the deployment workflow manually:

1. Go to **Actions** tab in your repository
2. Select **Deploy Documentation** workflow
3. Click **Run workflow**
4. Enter a test version (e.g., `1.0.0`)
5. Run the workflow

## 6. Create Your First Release

When you're ready to deploy documentation for a release:

1. Create a new release on GitHub (e.g., `v1.0.0`)
2. The workflow will automatically:
   - Extract the version from the tag
   - Create versioned documentation
   - Deploy to GitHub Pages

## 7. Build and Verify

Before pushing, you can build locally to check for errors:

```bash
npm run build
```

This will create a `build/` directory with the static site. You can serve it locally:

```bash
npm run serve
```

## Notes

- The documentation source files are in `docs/website/docs/`
- The sidebar structure is defined in `sidebars.js`
- Configuration is in `docusaurus.config.js`
- Versioned docs will be created in `versioned_docs/` when you run `npm run version`

