#!/usr/bin/env node
const { spawnSync } = require("child_process");
const { join } = require("path");
const { existsSync } = require("fs");

const BINARY_NAME = process.platform === "win32" ? "mobi.exe" : "mobi";
const BINARY_PATH = join(__dirname, "..", BINARY_NAME);

if (!existsSync(BINARY_PATH)) {
  console.error("⬇️  Descargando CLI mobi...");
  require("../install")(BINARY_PATH);
}

const result = spawnSync(BINARY_PATH, process.argv.slice(2), {
  stdio: "inherit",
});

process.exit(result.status ?? 1);
