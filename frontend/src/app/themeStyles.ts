import {
  themePalettes,
  type ThemeMode,
  type ThemeState,
  type ThemeTone,
} from "@/app/store/themeSlice";

const defaultPaletteName = "shadcnNeutral";

export function cx(...classes: Array<false | null | string | undefined>) {
  return classes.filter(Boolean).join(" ");
}

export function getThemeTone(theme: ThemeMode | ThemeState): ThemeTone {
  const mode = typeof theme === "string" ? theme : theme.mode;
  const paletteName = typeof theme === "string" ? defaultPaletteName : theme.paletteName;
  return themePalettes[paletteName]?.modes[mode] ?? themePalettes[defaultPaletteName].modes[mode];
}

export function getThemeButtonClassName(theme: ThemeMode | ThemeState) {
  const tone = getThemeTone(theme);

  return cx(
    "rounded-md border shadow-sm transition-colors",
    "disabled:pointer-events-none disabled:opacity-50",
    tone.border,
    tone.focusRing,
    tone.secondary,
    tone.secondaryForeground,
    tone.secondaryHover,
  );
}

export function getThemePrimaryButtonClassName(theme: ThemeMode | ThemeState) {
  const tone = getThemeTone(theme);

  return cx(
    "rounded-md shadow-sm transition-colors",
    "disabled:pointer-events-none disabled:opacity-50",
    tone.focusRing,
    tone.primary,
    tone.primaryForeground,
    tone.primaryHover,
  );
}
