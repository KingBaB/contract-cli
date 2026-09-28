"use strict";

const assert = require("node:assert/strict");
const crypto = require("node:crypto");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const { verifyArchiveChecksum } = require("./install.js");

test("verifyArchiveChecksum accepts the published SHA-256", () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), "contract-cli-install-test-"));
  try {
    const archive = path.join(dir, "contract-cli-1.2.0-linux-amd64.tar.gz");
    fs.writeFileSync(archive, "trusted archive");
    const digest = crypto.createHash("sha256").update("trusted archive").digest("hex");
    const checksums = path.join(dir, "checksums.txt");
    fs.writeFileSync(checksums, `${digest}  ${path.basename(archive)}\n`);

    verifyArchiveChecksum(archive, checksums);
  } finally {
    fs.rmSync(dir, { recursive: true, force: true });
  }
});

test("verifyArchiveChecksum rejects a mismatched archive", () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), "contract-cli-install-test-"));
  try {
    const archive = path.join(dir, "contract-cli-1.2.0-linux-amd64.tar.gz");
    fs.writeFileSync(archive, "tampered archive");
    const checksums = path.join(dir, "checksums.txt");
    fs.writeFileSync(checksums, `${"0".repeat(64)}  ${path.basename(archive)}\n`);

    assert.throws(() => verifyArchiveChecksum(archive, checksums), /checksum mismatch/);
  } finally {
    fs.rmSync(dir, { recursive: true, force: true });
  }
});
