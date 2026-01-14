#!/usr/bin/env node

const https = require('https');
const http = require('http');
const fs = require('fs');
const path = require('path');
const { pipeline } = require('stream');
const { promisify } = require('util');

const streamPipeline = promisify(pipeline);

// Read package version
const packageJson = require('../package.json');
const version = packageJson.version;

// Detect platform and architecture
const platform = process.platform;
const arch = process.arch;

// Map Node.js platform/arch to Go GOOS/GOARCH
const platformMap = {
  'win32': 'windows',
  'darwin': 'darwin',
  'linux': 'linux'
};

const archMap = {
  'x64': 'amd64',
  'arm64': 'arm64'
};

const goPlatform = platformMap[platform];
const goArch = archMap[arch];

if (!goPlatform || !goArch) {
  console.error(`Unsupported platform/architecture: ${platform}/${arch}`);
  process.exit(1);
}

// Construct binary name and download URL
const ext = platform === 'win32' ? '.exe' : '';
const binaryName = `codeloc-${goPlatform}-${goArch}${ext}`;
const downloadURL = `https://github.com/stripsior/loc.cli/releases/download/v${version}/${binaryName}`;

// Create cache directory
const cacheDir = path.join(__dirname, '..', '.bin-cache');
if (!fs.existsSync(cacheDir)) {
  fs.mkdirSync(cacheDir, { recursive: true });
}

const binaryPath = path.join(cacheDir, binaryName);

// Download binary
async function downloadBinary() {
  console.log(`Downloading CodeLoc binary for ${goPlatform}/${goArch}...`);
  console.log(`URL: ${downloadURL}`);

  return new Promise((resolve, reject) => {
    const protocol = downloadURL.startsWith('https') ? https : http;

    protocol.get(downloadURL, (response) => {
      // Handle redirects
      if (response.statusCode === 302 || response.statusCode === 301) {
        const redirectURL = response.headers.location;
        console.log(`Following redirect to: ${redirectURL}`);

        const redirectProtocol = redirectURL.startsWith('https') ? https : http;
        redirectProtocol.get(redirectURL, async (redirectResponse) => {
          if (redirectResponse.statusCode !== 200) {
            reject(new Error(`Failed to download binary: HTTP ${redirectResponse.statusCode}`));
            return;
          }

          try {
            await streamPipeline(redirectResponse, fs.createWriteStream(binaryPath));
            resolve();
          } catch (err) {
            reject(err);
          }
        }).on('error', reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download binary: HTTP ${response.statusCode}`));
        return;
      }

      streamPipeline(response, fs.createWriteStream(binaryPath))
        .then(resolve)
        .catch(reject);
    }).on('error', reject);
  });
}

// Make binary executable (Unix only)
function makeExecutable() {
  if (platform !== 'win32') {
    fs.chmodSync(binaryPath, 0o755);
  }
}

// Main installation process
(async () => {
  try {
    // Check if binary already exists
    if (fs.existsSync(binaryPath)) {
      console.log('Binary already exists, skipping download.');
      makeExecutable();
      console.log('CodeLoc installation complete!');
      return;
    }

    await downloadBinary();
    makeExecutable();

    console.log('CodeLoc binary downloaded successfully!');
    console.log('You can now use: npx codeloc');
  } catch (error) {
    console.error('Installation failed:', error.message);
    console.error('\nYou can try:');
    console.error('1. Check your internet connection');
    console.error('2. Verify the release exists on GitHub');
    console.error('3. Download the binary manually from:');
    console.error(`   ${downloadURL}`);
    process.exit(1);
  }
})();
