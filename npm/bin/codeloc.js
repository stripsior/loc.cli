#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

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

// Construct binary name
const ext = platform === 'win32' ? '.exe' : '';
const binaryName = `codeloc-${goPlatform}-${goArch}${ext}`;

// Find binary path
const binaryPath = path.join(__dirname, '..', '.bin-cache', binaryName);

if (!fs.existsSync(binaryPath)) {
  console.error(`Binary not found: ${binaryPath}`);
  console.error('Please try reinstalling the package: npm install codeloc');
  process.exit(1);
}

// Execute binary with forwarded arguments
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  shell: false
});

child.on('error', (error) => {
  console.error(`Failed to execute binary: ${error.message}`);
  process.exit(1);
});

child.on('exit', (code) => {
  process.exit(code || 0);
});
