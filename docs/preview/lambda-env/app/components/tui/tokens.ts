// Design tokens synced with Go source of truth:
// src/lambda-env/internal/tui/theme/colors.go
export const C = {
  // Backgrounds
  bg:      '#0D0D0D',
  surface: '#1E1E2E',
  surface2:'#252540',

  // Accent / Primary
  accent:  '#8B6AF4',
  accentDim:'#7D56F440',
  accentBorder:'#7D56F460',

  // Semantic
  success: '#04B575',
  successDim: '#04B57530',
  error:   '#FF4672',
  errorDim:'#FF467220',
  warn:    '#F4D03F',
  warnDim: '#F4D03F20',

  // Text
  textPrimary:   '#FFFFFF',
  textSecondary: '#888888',
  textMuted:     '#444466',
  dimmed:        '#909090', // refined from old #626262; passes WCAG AA on #1A1A1A

  // Borders
  border:    '#2A2A4A',
  borderFocus: '#8B6AF4',
};
