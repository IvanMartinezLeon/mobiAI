/**
 * MOBI AI Custom Header Extension for Pi Coding Agent
 */

import type { ExtensionAPI, Theme } from "@earendil-works/pi-coding-agent";
import { VERSION } from "@earendil-works/pi-coding-agent";

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
	// Set custom header immediately on load (if UI is available)
	pi.on("session_start", async (_event, ctx) => {
		if (ctx.hasUI) {
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
		}
	});

	// Command to restore built-in header
	pi.registerCommand("builtin-header", {
		description: "Restore built-in header with keybinding hints",
		handler: async (_args, ctx) => {
			ctx.ui.setHeader(undefined);
			ctx.ui.notify("Built-in header restored", "info");
		},
	});
}
