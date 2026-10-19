/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { describe, expect, it } from 'vitest'

// DESIGN.md ở gốc kho mã chép lại token của theme.css cho trợ lý AI; kiểm thử giữ hai nơi khớp nhau.
const css = readFileSync(`${import.meta.dirname}/theme.css`, 'utf8')
const design = readFileSync(`${import.meta.dirname}/../../../../DESIGN.md`, 'utf8')

function cssTokens(selector: string): Record<string, string> {
  const start = css.indexOf(`${selector} {`)
  const block = css.slice(start, css.indexOf('}', start))
  return Object.fromEntries([...block.matchAll(/--([\w-]+):\s*(#[0-9a-f]{6});/gi)].map((m) => [m[1]!, m[2]!.toLowerCase()]))
}

function designColors(): Record<string, string> {
  const front = design.split(/^---$/m)[1] ?? ''
  const block = front.slice(front.indexOf('colors:'), front.indexOf('typography:'))
  return Object.fromEntries([...block.matchAll(/^\s+([\w-]+):\s*"(#[0-9a-f]{6})"/gim)].map((m) => [m[1]!, m[2]!.toLowerCase()]))
}

describe('DESIGN.md khớp theme.css (W12)', () => {
  const colors = designColors()
  const light = cssTokens(':root')
  const dark = cssTokens('.dark')

  it('mọi màu chế độ sáng của theme.css có trong DESIGN.md với cùng giá trị', () => {
    expect(Object.keys(light).length).toBeGreaterThan(0)
    for (const [name, hex] of Object.entries(light)) expect(colors[name], name).toBe(hex)
  })

  it('mọi màu chế độ tối có trong DESIGN.md dưới tên <token>-dark', () => {
    for (const [name, hex] of Object.entries(dark)) expect(colors[`${name}-dark`], `${name}-dark`).toBe(hex)
  })

  it('DESIGN.md không có màu không tồn tại trong theme.css', () => {
    for (const name of Object.keys(colors)) {
      const base = name.endsWith('-dark') ? name.slice(0, -5) : name
      expect(name.endsWith('-dark') ? dark[base] : light[base], name).toBeDefined()
    }
  })
})
