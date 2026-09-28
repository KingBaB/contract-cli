"use strict";

const fs = require("fs");
const crypto = require("crypto");
const os = require("os");
const path = require("path");
const { execFileSync, spawnSync } = require("child_process");

const pkg = require("../package.json");
const binaryName = (pkg.config && pkg.config.binaryName) || "contract-cli";
const version = pkg.version;
const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};
const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

const platform = platformMap[process.platform];
const arch = archMap[process.arch];
if (!platform || !arch) {
  console.error(`Unsupported platform: ${process.platform}-${process.arch}`);
  process.exit(1);
}

const isWindows = process.platform === "win32";
const archiveExt = isWindows ? ".zip" : ".tar.gz";
const archiveName = `${binaryName}-${version}-${platform}-${arch}${archiveExt}`;
const rootDir = path.join(__dirname, "..");
const binDir = path.join(rootDir, "bin");
const destination = path.join(binDir, binaryName + (isWindows ? ".exe" : ""));

function resolveDownloadBaseURL() {
  const template =
    process.env.CONTRACT_CLI_DOWNLOAD_BASE_URL_TEMPLATE ||
    process.env.npm_package_config_downloadBaseURLTemplate ||
    (pkg.config && pkg.config.downloadBaseURLTemplate) ||
    "";
  if (!template) {
    return "";
  }
  return template.replace(/\{version\}/g, version).replace(/\/$/, "");
}

function commandExists(command) {
  const checker = process.platform === "win32" ? "where" : "which";
  const args = [command];
  const result = spawnSync(checker, args, { stdio: "ignore", shell: false });
  return result.status === 0;
}

function downloadArchive(downloadURL, archivePath) {
  if (!commandExists("curl")) {
    throw new Error("curl not found");
  }

  const curlArgs = [
    "--fail",
    "--location",
    "--silent",
    "--show-error",
    "--connect-timeout",
    "10",
    "--max-time",
    "120",
    "--output",
    archivePath,
    downloadURL,
  ];
  if (isWindows) {
    curlArgs.unshift("--ssl-revoke-best-effort");
  }
  execFileSync("curl", curlArgs, { stdio: ["ignore", "ignore", "pipe"] });
}

function verifyArchiveChecksum(archivePath, checksumsPath) {
  const archiveNameToVerify = path.basename(archivePath);
  const line = fs
    .readFileSync(checksumsPath, "utf8")
    .split(/\r?\n/)
    .map((value) => value.trim())
    .find((value) => value.endsWith(`  ${archiveNameToVerify}`) || value.endsWith(` *${archiveNameToVerify}`));
  if (!line) {
    throw new Error(`checksum for ${archiveNameToVerify} not found`);
  }
  const expected = line.split(/\s+/)[0].toLowerCase();
  if (!/^[a-f0-9]{64}$/.test(expected)) {
    throw new Error(`invalid checksum for ${archiveNameToVerify}`);
  }
  const actual = crypto.createHash("sha256").update(fs.readFileSync(archivePath)).digest("hex");
  if (actual !== expected) {
    throw new Error(`checksum mismatch for ${archiveNameToVerify}`);
  }
}

function extractArchive(archivePath, tempDir) {
  if (isWindows) {
    execFileSync(
      "powershell",
      [
        "-Command",
        `Expand-Archive -Path '${archivePath}' -DestinationPath '${tempDir}' -Force`,
      ],
      { stdio: "ignore" }
    );
    return;
  }
  execFileSync("tar", ["-xzf", archivePath, "-C", tempDir], { stdio: "ignore" });
}

function installFromArchive(archivePath) {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "contract-cli-"));

  try {
    extractArchive(archivePath, tempDir);
    const extractedBinary = path.join(tempDir, binaryName + (isWindows ? ".exe" : ""));
    if (!fs.existsSync(extractedBinary)) {
      throw new Error(`binary ${path.basename(extractedBinary)} not found in archive`);
    }
    fs.mkdirSync(binDir, { recursive: true });
    fs.copyFileSync(extractedBinary, destination);
    fs.chmodSync(destination, 0o755);
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

function installFromBundledArchive() {
  const archivePath = path.join(rootDir, "dist", "release-assets", archiveName);
  if (!fs.existsSync(archivePath)) {
    return false;
  }

  verifyArchiveChecksum(archivePath, path.join(rootDir, "dist", "release-assets", "checksums.txt"));
  installFromArchive(archivePath);
  return true;
}

function installFromDownload(downloadBaseURL) {
  if (!downloadBaseURL) {
    throw new Error("download base URL template not configured");
  }

  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "contract-cli-"));
  const archivePath = path.join(tempDir, archiveName);
  const checksumsPath = path.join(tempDir, "checksums.txt");
  const downloadURL = `${downloadBaseURL}/${archiveName}`;

  try {
    downloadArchive(downloadURL, archivePath);
    downloadArchive(`${downloadBaseURL}/checksums.txt`, checksumsPath);
    verifyArchiveChecksum(archivePath, checksumsPath);
    installFromArchive(archivePath);
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

function install() {
  const downloadBaseURL = resolveDownloadBaseURL();

  if (installFromBundledArchive()) {
    console.log(`${binaryName} ${version} installed from bundled npm assets`);
    return;
  }

  if (downloadBaseURL) {
    installFromDownload(downloadBaseURL);
    console.log(`${binaryName} ${version} installed from release assets`);
    return;
  }

  throw new Error("download base URL template not configured and no bundled release archive is available");
}

if (require.main === module) {
  try {
    install();
  } catch (error) {
    console.error(`Failed to install ${binaryName}: ${error.message}`);
    console.error("Verify the release platform, download URL, and published checksums.txt, then retry installation.");
    process.exit(1);
  }
}

module.exports = { verifyArchiveChecksum };
