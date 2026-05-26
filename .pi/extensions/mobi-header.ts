/**
 * MOBI AI Custom Header Extension for Pi Coding Agent
 */

import type { ExtensionAPI, ExtensionContext, Theme } from "@earendil-works/pi-coding-agent";
import { VERSION } from "@earendil-works/pi-coding-agent";
import { readFileSync, existsSync } from "fs";
import { join } from "path";

function detectFramework(cwd: string): string {
	try {
		if (existsSync(join(cwd, "pubspec.yaml"))) {
			const content = readFileSync(join(cwd, "pubspec.yaml"), "utf-8");
			const name = content.match(/^name:\s*(.+)$/m)?.[1]?.trim() || "";
			const isFlutter = content.includes("flutter:");
			return name ? `${name} (${isFlutter ? "Flutter" : "Dart"})` : isFlutter ? "Flutter" : "Dart";
		}
		if (existsSync(join(cwd, "package.json"))) {
			const pkg = JSON.parse(readFileSync(join(cwd, "package.json"), "utf-8"));
			let label = pkg.name || "Node.js";
			const deps = Object.keys(pkg.dependencies || {});
			if (deps.includes("next")) label += " + Next.js";
			else if (deps.includes("react")) label += " + React";
			else if (deps.includes("vue")) label += " + Vue";
			return `${label} (Node.js)`;
		}
		if (existsSync(join(cwd, "Cargo.toml"))) {
			const content = readFileSync(join(cwd, "Cargo.toml"), "utf-8");
			const name = content.match(/^name\s*=\s*"(.+)"$/m)?.[1] || "";
			return name ? `${name} (Rust)` : "Rust";
		}
		if (existsSync(join(cwd, "go.mod"))) {
			const content = readFileSync(join(cwd, "go.mod"), "utf-8");
			const name = content.match(/^module\s+(.+)$/m)?.[1] || "";
			return name ? `${name} (Go)` : "Go";
		}
		if (existsSync(join(cwd, "pyproject.toml"))) {
			const content = readFileSync(join(cwd, "pyproject.toml"), "utf-8");
			const name = content.match(/^name\s*=\s*"(.+)"$/m)?.[1] || "";
			return name ? `${name} (Python)` : "Python";
		}
		if (existsSync(join(cwd, "pom.xml"))) {
			const content = readFileSync(join(cwd, "pom.xml"), "utf-8");
			const name = content.match(/<artifactId>(.+?)<\/artifactId>/)?.[1] || "";
			return name ? `${name} (Java/Maven)` : "Java/Maven";
		}
		if (existsSync(join(cwd, "build.gradle")) || existsSync(join(cwd, "build.gradle.kts"))) {
			return "Java/Kotlin (Gradle)";
		}
		if (existsSync(join(cwd, "composer.json"))) {
			const pkg = JSON.parse(readFileSync(join(cwd, "composer.json"), "utf-8"));
			return `${pkg.name || "PHP"} (PHP)`;
		}
		if (existsSync(join(cwd, "mix.exs"))) return "Elixir (Phoenix)";
		if (existsSync(join(cwd, "Deno.json")) || existsSync(join(cwd, "deno.jsonc"))) return "Deno";
		return "";
	} catch {
		return "";
	}
}

function formatNum(n: number): string {
	return n >= 1000 ? `${(n / 1000).toFixed(1)}k` : `${n}`;
}

function stripAnsi(str: string): string {
	return str.replace(/\x1B\[[0-9;]*[a-zA-Z]/g, "");
}

function getMobiLogo(theme: Theme): string[] {
	const cyan = (text: string) => theme.fg("accent", text);
	const logo = [
		"",
		" ███╗   ███╗ ██████╗ ██████╗ ██╗     █████╗ ██╗",
		" ████╗ ████║██╔═══██╗██╔══██╗██║    ██╔══██╗██║",
		" ██╔████╔██║██║   ██║██████╔╝██║    ███████║██║",
		" ██║╚██╔╝██║██║   ██║██╔══██╗██║    ██╔══██║██║",
		" ██║ ╚═╝ ██║╚██████╔╝██████╔╝██║    ██║  ██║██║",
		" ╚═╝     ╚═╝ ╚═════╝ ╚═════╝ ╚═╝    ╚═╝  ╚═╝╚═╝",
		""
	];
	return logo.map(line => cyan(line));
}

export default function (pi: ExtensionAPI) {
	let cachedFramework: string | null = null;
	let sessionCtx: ExtensionContext | null = null;
	let totalInput = 0;
	let totalOutput = 0;

	pi.on("message_end", (_event, _ctx) => {
		const msg = _event.message as Record<string, unknown>;
		const usage = msg?.usage as Record<string, unknown> | undefined;
		if (usage && typeof usage.input === "number" && typeof usage.output === "number") {
			totalInput += usage.input as number;
			totalOutput += usage.output as number;
		}
	});

	pi.on("session_start", async (_event, ctx) => {
		if (ctx.hasUI) {
			if (cachedFramework === null) {
				cachedFramework = detectFramework(ctx.cwd);
			}
			sessionCtx = ctx;
			const framework = cachedFramework;

			ctx.ui.setHeader((_tui, theme) => {
				return {
					render(_width: number): string[] {
						const logoLines = getMobiLogo(theme);
						const subtitle = `${theme.fg("muted", "   MOBI AI Custom Agent Interface")}${theme.fg("dim", ` v${VERSION}`)}`;
						return [...logoLines, subtitle, ""];
					},
					invalidate() {},
				};
			});

			ctx.ui.setFooter((_tui, theme, footerData) => {
				const branch = footerData.getGitBranch();
				return {
					render(width: number): string[] {
						const left: string[] = [];
						if (framework) left.push(`${theme.fg("muted", "Framework:")} ${theme.fg("success", framework)}`);
						if (branch) left.push(`${theme.fg("muted", "Branch:")} ${theme.fg("warning", branch)}`);

						const modelName = sessionCtx?.model?.id ?? "";
						const usage = sessionCtx?.getContextUsage();
						const pct = usage?.percent != null ? usage.percent : null;

						const right: string[] = [];
						if (modelName) right.push(`${theme.fg("muted", "Model:")} ${theme.fg("accent", modelName)}`);
						if (pct != null) {
							const color = pct > 80 ? "error" : pct > 50 ? "warning" : "success";
							right.push(`${theme.fg("muted", "Context:")} ${theme.fg(color, `${pct.toFixed(0)}%`)}`);
						}
						if (totalInput + totalOutput > 0) {
							right.push(`${theme.fg("muted", "Tokens:")} ${theme.fg("accent", `${formatNum(totalInput)}↑`)} ${theme.fg("accent", `${formatNum(totalOutput)}↓`)}`);
						}

						const leftStr = left.join(` ${theme.fg("dim", "|")} `);
						const rightStr = right.join(` ${theme.fg("dim", "|")} `);
						const sep = leftStr && rightStr ? " ".repeat(Math.max(2, width - stripAnsi(leftStr).length - stripAnsi(rightStr).length - 2)) : "";
						const line = [leftStr, rightStr].filter(Boolean).join(sep);

						return line ? ["", line, ""] : [""];
					},
					invalidate() {},
				};
			});
		}
	});

	pi.registerCommand("builtin-header", {
		description: "Restore built-in header with keybinding hints",
		handler: async (_args, ctx) => {
			ctx.ui.setHeader(undefined);
			ctx.ui.setFooter(undefined);
			ctx.ui.notify("Built-in header restored", "info");
		},
	});
}
