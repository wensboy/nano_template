import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type ThemeMode = "dark" | "light";
export type ThemePaletteName = "shadcnNeutral";

export type ThemeTone = {
  border: string;
  borderStrong: string;
  destructive: string;
  destructiveForeground: string;
  destructiveHover: string;
  fileDone: string;
  focusRing: string;
  foreground: string;
  iconMuted: string;
  muted: string;
  mutedForeground: string;
  overlay: string;
  page: string;
  primary: string;
  primaryForeground: string;
  primaryHover: string;
  progress: string;
  progressTrack: string;
  secondary: string;
  secondaryForeground: string;
  secondaryHover: string;
  shadow: string;
  subtle: string;
  subtleForeground: string;
  surface: string;
  surfaceForeground: string;
};

export type ThemePaletteConfig = {
  label: string;
  modes: Record<ThemeMode, ThemeTone>;
};

export type ThemeState = {
  mode: ThemeMode;
  paletteName: ThemePaletteName;
};

const initialState: ThemeState = {
  mode: "dark",
  paletteName: "shadcnNeutral",
};

export const themePalettes: Record<ThemePaletteName, ThemePaletteConfig> = {
  shadcnNeutral: {
    label: "Shadcn Neutral",
    modes: {
      light: {
        border: "border-zinc-200",
        borderStrong: "border-zinc-300",
        destructive: "bg-red-500",
        destructiveForeground: "text-white",
        destructiveHover: "hover:bg-red-600",
        fileDone: "border-emerald-500 bg-emerald-600 text-white",
        focusRing: "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-400 focus-visible:ring-offset-2 focus-visible:ring-offset-white",
        foreground: "text-zinc-950",
        iconMuted: "text-zinc-500",
        muted: "bg-zinc-100",
        mutedForeground: "text-zinc-600",
        overlay: "bg-zinc-950/45",
        page: "bg-zinc-50",
        primary: "bg-zinc-950",
        primaryForeground: "text-zinc-50",
        primaryHover: "hover:bg-zinc-800",
        progress: "text-emerald-600",
        progressTrack: "text-zinc-200",
        secondary: "bg-white",
        secondaryForeground: "text-zinc-950",
        secondaryHover: "hover:bg-zinc-100",
        shadow: "shadow-2xl",
        subtle: "bg-zinc-50",
        subtleForeground: "text-zinc-500",
        surface: "bg-white",
        surfaceForeground: "text-zinc-950",
      },
      dark: {
        border: "border-zinc-800",
        borderStrong: "border-zinc-700",
        destructive: "bg-red-600",
        destructiveForeground: "text-white",
        destructiveHover: "hover:bg-red-500",
        fileDone: "border-emerald-500 bg-emerald-600 text-white",
        focusRing: "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-500 focus-visible:ring-offset-2 focus-visible:ring-offset-zinc-950",
        foreground: "text-zinc-50",
        iconMuted: "text-zinc-400",
        muted: "bg-zinc-800",
        mutedForeground: "text-zinc-300",
        overlay: "bg-zinc-950/70",
        page: "bg-zinc-950",
        primary: "bg-zinc-50",
        primaryForeground: "text-zinc-950",
        primaryHover: "hover:bg-zinc-200",
        progress: "text-emerald-500",
        progressTrack: "text-zinc-700",
        secondary: "bg-zinc-900",
        secondaryForeground: "text-zinc-50",
        secondaryHover: "hover:bg-zinc-800",
        shadow: "shadow-2xl",
        subtle: "bg-zinc-900",
        subtleForeground: "text-zinc-400",
        surface: "bg-zinc-900",
        surfaceForeground: "text-zinc-50",
      },
    },
  },
};

const themeSlice = createSlice({
  name: "theme",
  initialState,
  reducers: {
    setThemeMode(state, action: PayloadAction<ThemeMode>) {
      state.mode = action.payload;
    },
    setThemePaletteName(state, action: PayloadAction<ThemePaletteName>) {
      state.paletteName = action.payload;
    },
    toggleThemeMode(state) {
      state.mode = state.mode === "dark" ? "light" : "dark";
    },
  },
});

export const { setThemeMode, setThemePaletteName, toggleThemeMode } = themeSlice.actions;
export default themeSlice.reducer;
