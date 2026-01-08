#!/usr/bin/env node
/**
 * Script to sync Italian translations for versioned docs.
 * This script creates version.json files for each version in the Italian i18n directory.
 */

const fs = require('fs');
const path = require('path');

const versionsJsonPath = path.join(__dirname, '..', 'versions.json');
const versionedDocsPath = path.join(__dirname, '..', 'versioned_docs');
const i18nItBasePath = path.join(__dirname, '..', 'i18n', 'it', 'docusaurus-plugin-content-docs');

// Read versions.json
let versions = [];
if (fs.existsSync(versionsJsonPath)) {
  try {
    const versionsContent = fs.readFileSync(versionsJsonPath, 'utf8');
    versions = JSON.parse(versionsContent);
  } catch (error) {
    console.error('Error reading versions.json:', error);
    process.exit(1);
  }
}

// Also check for versioned_docs directories if versions.json is empty or incomplete
if (versions.length === 0 && fs.existsSync(versionedDocsPath)) {
  try {
    const entries = fs.readdirSync(versionedDocsPath, { withFileTypes: true });
    versions = entries
      .filter(entry => entry.isDirectory() && entry.name.startsWith('version-'))
      .map(entry => entry.name.replace('version-', ''));
    console.log(`Found ${versions.length} versions from versioned_docs directory: ${versions.join(', ')}`);
  } catch (error) {
    console.warn(`Warning: Could not read versioned_docs directory: ${error.message}`);
  }
}

// Ensure i18n/it/docusaurus-plugin-content-docs directory exists
if (!fs.existsSync(i18nItBasePath)) {
  fs.mkdirSync(i18nItBasePath, { recursive: true });
}

// Process each version
versions.forEach((version) => {
  const versionDir = path.join(i18nItBasePath, `version-${version}`);
  const versionJsonPath = path.join(versionDir, 'version.json');
  
  // Create version directory if it doesn't exist
  if (!fs.existsSync(versionDir)) {
    fs.mkdirSync(versionDir, { recursive: true });
  }
  
  // Create or update version.json with the version label
  const versionJson = {
    'version.label': {
      message: version,
      description: `The label for version ${version} in the version dropdown`
    }
  };
  
  // If version.json exists, merge with existing content
  let existingJson = {};
  if (fs.existsSync(versionJsonPath)) {
    try {
      const existingContent = fs.readFileSync(versionJsonPath, 'utf8');
      existingJson = JSON.parse(existingContent);
    } catch (error) {
      console.warn(`Warning: Could not parse existing ${versionJsonPath}, overwriting`);
    }
  }
  
  // Merge version label with existing content
  const mergedJson = { ...existingJson, ...versionJson };
  
  // Write version.json
  fs.writeFileSync(versionJsonPath, JSON.stringify(mergedJson, null, 2) + '\n');
  console.log(`✓ Created/updated ${versionJsonPath}`);
});

// Also ensure current version has the correct label
const currentJsonPath = path.join(i18nItBasePath, 'current.json');
if (fs.existsSync(currentJsonPath)) {
  try {
    const currentContent = fs.readFileSync(currentJsonPath, 'utf8');
    const currentJson = JSON.parse(currentContent);
    
    // Ensure version.label exists for current
    if (!currentJson['version.label']) {
      currentJson['version.label'] = {
        message: 'Next 🚧',
        description: 'The label for the current version in the version dropdown'
      };
      fs.writeFileSync(currentJsonPath, JSON.stringify(currentJson, null, 2) + '\n');
      console.log('✓ Updated current.json with version.label');
    }
  } catch (error) {
    console.warn(`Warning: Could not update current.json: ${error.message}`);
  }
}

console.log('✓ Version translations synced successfully');

