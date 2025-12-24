# Documentation PR Checks

This document explains the automated checks that run on pull requests that modify documentation.

## Automated Checks

When you open a PR that modifies documentation files, the following checks run automatically:

### Documentation PR Check

This workflow (`docs-pr-check.yml`) validates that:

1. ✅ Documentation builds successfully
2. ✅ No broken markdown links
3. ✅ No syntax errors
4. ✅ All dependencies are correctly installed

## What Triggers the Check

The check runs automatically when your PR modifies:
- Files in `docs/website/`
- Documentation workflow files (`.github/workflows/docs-*.yml`)
- Markdown files in `docs/`

## Required Status Checks

To ensure PRs cannot be merged until documentation builds successfully:

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Branches**
3. Edit or create a branch protection rule for your main branch
4. Under **Require status checks to pass before merging**, enable:
   - ✅ **Documentation PR Check**

## Common Issues

### Build Fails Due to Broken Links

**Error**: `Docusaurus found broken links!`

**Solution**: 
- Check the build logs for the list of broken links
- Fix all broken markdown links in your documentation
- Ensure all referenced files exist

### Build Fails Due to Syntax Errors

**Error**: `MDX compilation failed`

**Solution**:
- Check for unclosed code blocks (```)
- Ensure angle brackets in code blocks are properly formatted
- Verify markdown syntax is correct

### Node.js Version Issues

**Error**: `Minimum Node.js version not met`

**Solution**:
- Ensure you're using Node.js 20+ locally
- The CI uses Node.js 20 automatically

## Local Testing

Before pushing, test your documentation locally:

```bash
cd docs/website
npm install
npm run build
```

This will catch most issues before opening a PR.

