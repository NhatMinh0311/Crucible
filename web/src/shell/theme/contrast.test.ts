/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { describe, expect, it } from 'vitest'

// Vitest bỏ nội dung CSS khi import, nên đọc thẳng tệp.
const css = readFileSync(`${import.meta.dirname}/theme.css`, 'utf8')

/** Đọc token hex trong một khối `selector { ... }` của theme.css. */
function tokens(selector: string): Record<string, string> {
  const start = css.indexOf(`${selector} {`)
  const block = css.slice(start, css.indexOf('}', start))
  return Object.fromEntries([...block.matchAll(/--([\w-]+):\s*(#[0-9a-f]{6});/gi)].map((m) => [m[1]!, m[2]!]))
}

function luminance(hex: string) {
  const [r, g, b] = [1, 3, 5].map((i) => {
    const c = parseInt(hex.slice(i, i + 2), 16) / 255
    return c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4
  })
  return 0.2126 * r! + 0.7152 * g! + 0.0722 * b!
}

export function contrast(a: string, b: string) {
  const [hi, lo] = [luminance(a), luminance(b)].sort((x, y) => y - x)
  return (hi! + 0.05) / (lo! + 0.05)
}

// Cặp chữ trên nền xuất hiện trong khung chung và thành phần dùng chung.
const PAIRS: [string, string][] = [
  ['foreground', 'background'],
  ['foreground', 'surface'],
  ['foreground', 'accent'],
  ['muted-foreground', 'background'],
  ['muted-foreground', 'surface'],
  ['primary', 'background'],
  ['primary', 'surface'],
  ['primary', 'accent'],
  ['primary-foreground', 'primary'],
  ['signal', 'background'],
  ['signal', 'surface'],
  ['signal', 'signal-soft'],
  ['signal-foreground', 'signal'],
  ['success', 'success-soft'],
  ['danger', 'danger-soft'],
  ['danger', 'surface'],
  ['info', 'info-soft'],
  ['muted-foreground', 'muted'],
]

describe('độ tương phản của token (UX7, W12)', () => {
  for (const [mode, selector] of [['sáng', ':root'], ['tối', '.dark']] as const) {
    const t = { ...tokens(':root'), ...tokens(selector) }
    for (const [fg, bg] of PAIRS) {
      it(`${mode}: ${fg} trên ${bg} đạt 4,5:1`, () => {
        expect(t[fg], fg).toBeDefined()
        expect(t[bg], bg).toBeDefined()
        expect(contrast(t[fg]!, t[bg]!)).toBeGreaterThanOrEqual(4.5)
      })
    }
  }

  it('terminal: chữ trên nền tối đạt 4,5:1', () => {
    const t = tokens(':root')
    expect(contrast(t['terminal-foreground']!, t.terminal!)).toBeGreaterThanOrEqual(4.5)
  })
})
