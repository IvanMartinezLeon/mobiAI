const { createWriteStream, unlinkSync, renameSync } = require("fs");
const { get } = require("https");
const { platform, arch } = process;

const REPO = "IvanMartinezLeon/mobiAI";

function getOS() {
  if (platform === "darwin") return "darwin";
  if (platform === "linux") return "linux";
  if (platform === "win32") return "windows";
  throw new Error(`Platform not supported: ${platform}`);
}

function getArch() {
  if (arch === "x64") return "amd64";
  if (arch === "arm64") return "arm64";
  throw new Error(`Architecture not supported: ${arch}`);
}

module.exports = function download(destPath) {
  return new Promise((resolve, reject) => {
    const url = `https://github.com/${REPO}/releases/latest/download/mobi_${getOS()}_${getArch()}.tar.gz`;

    get("https://api.github.com/repos/" + REPO + "/releases/latest", {
      headers: { "User-Agent": "mobi-cli" },
    }, (res) => {
      let body = "";
      res.on("data", (chunk) => body += chunk);
      res.on("end", () => {
        const tag = JSON.parse(body).tag_name;
        const binUrl = `https://github.com/${REPO}/releases/download/${tag}/mobi_${getOS()}_${getArch()}.tar.gz`;
        downloadAndExtract(binUrl, destPath).then(resolve).catch(reject);
      });
    }).on("error", reject);
  });
};

function downloadAndExtract(url, destPath) {
  return new Promise((resolve, reject) => {
    const tmp = destPath + ".tmp";
    const gunzip = require("child_process").spawn("tar", ["-xzf", "-", "-C", require("path").dirname(destPath)]);
    const file = createWriteStream(tmp, { mode: 0o755 });

    get(url, { headers: { "User-Agent": "mobi-cli" } }, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        downloadAndExtract(res.headers.location, destPath).then(resolve).catch(reject);
        return;
      }
      res.pipe(gunzip.stdin);
      res.on("error", reject);
    }).on("error", reject);

    gunzip.stdout.pipe(file);
    file.on("finish", () => {
      file.close();
      try { unlinkSync(destPath); } catch {}
      renameSync(tmp, destPath);
      resolve();
    });
    gunzip.on("error", reject);
  });
}
