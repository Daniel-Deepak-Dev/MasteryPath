/**
 * Theme constants for MasteryPath
 * Extended with dashboard-specific colors, gradients, and mastery level palette.
 */
import { Platform } from 'react-native';

const tintColorLight = '#6C63FF';
const tintColorDark = '#A78BFA';

export const Colors = {
  light: {
    text: '#1A1A2E',
    textSecondary: '#64748B',
    background: '#F8FAFC',
    surface: '#FFFFFF',
    tint: tintColorLight,
    icon: '#64748B',
    tabIconDefault: '#94A3B8',
    tabIconSelected: tintColorLight,
    border: '#E2E8F0',
    // Dashboard specific
    cardGradientStart: '#6C63FF',
    cardGradientEnd: '#A78BFA',
    masteryBg: 'rgba(108, 99, 255, 0.08)',
    chipBg: 'rgba(108, 99, 255, 0.1)',
    chipText: '#6C63FF',
    headerGradientStart: '#6C63FF',
    headerGradientEnd: '#818CF8',
  },
  dark: {
    text: '#F1F5F9',
    textSecondary: '#94A3B8',
    background: '#0F172A',
    surface: '#1E293B',
    tint: tintColorDark,
    icon: '#94A3B8',
    tabIconDefault: '#64748B',
    tabIconSelected: tintColorDark,
    border: '#334155',
    // Dashboard specific
    cardGradientStart: '#4F46E5',
    cardGradientEnd: '#7C3AED',
    masteryBg: 'rgba(167, 139, 250, 0.12)',
    chipBg: 'rgba(167, 139, 250, 0.15)',
    chipText: '#A78BFA',
    headerGradientStart: '#4F46E5',
    headerGradientEnd: '#6366F1',
  },
};

// Mastery level colors (0-100 scale)
export const MasteryColors = {
  beginner: '#F87171',    // red-400   (0-25%)
  developing: '#FB923C',  // orange-400 (25-50%)
  proficient: '#FBBF24',  // amber-400  (50-75%)
  expert: '#34D399',      // emerald-400 (75-100%)
};

export const getMasteryColor = (mastery: number): string => {
  if (mastery >= 75) return MasteryColors.expert;
  if (mastery >= 50) return MasteryColors.proficient;
  if (mastery >= 25) return MasteryColors.developing;
  return MasteryColors.beginner;
};

export const getMasteryLabel = (mastery: number): string => {
  if (mastery >= 75) return 'Expert';
  if (mastery >= 50) return 'Proficient';
  if (mastery >= 25) return 'Developing';
  return 'Beginner';
};

// Radar chart axis colors (cycling)
export const RadarAxisColors = [
  '#6C63FF', '#F472B6', '#34D399', '#FBBF24', '#60A5FA',
  '#A78BFA', '#FB923C', '#2DD4BF', '#E879F9', '#38BDF8',
];

export const Fonts = Platform.select({
  ios: {
    sans: 'system-ui',
    serif: 'ui-serif',
    rounded: 'ui-rounded',
    mono: 'ui-monospace',
  },
  default: {
    sans: 'normal',
    serif: 'serif',
    rounded: 'normal',
    mono: 'monospace',
  },
  web: {
    sans: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif",
    serif: "Georgia, 'Times New Roman', serif",
    rounded: "'SF Pro Rounded', 'Hiragino Maru Gothic ProN', Meiryo, 'MS PGothic', sans-serif",
    mono: "SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace",
  },
});
