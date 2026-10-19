import ts from 'typescript'
import { describe, expect, it } from 'vitest'
import { LANGUAGES, NAMESPACES, resources } from '@/shell/i18n/i18n'

type Tree = { [key: string]: string | Tree }

function flatten(tree: Tree, prefix = ''): Map<string, string> {
  const out = new Map<string, string>()
  for (const [k, v] of Object.entries(tree)) {
    const key = prefix ? `${prefix}.${k}` : k
    if (typeof v === 'string') out.set(key, v)
    else for (const [kk, vv] of flatten(v, key)) out.set(kk, vv)
  }
  return out
}

const variables = (text: string) => [...text.matchAll(/{{\s*(\w+)\s*}}/g)].map((m) => m[1]).sort()

describe('tệp dịch (AC-P2.1, W3)', () => {
  for (const ns of NAMESPACES) {
    it(`không gian ${ns}: vi và en có cùng tập khóa, cùng biến, không có bản dịch rỗng`, () => {
      const [vi, en] = LANGUAGES.map((l) => flatten(resources[l][ns] as Tree))
      expect([...en!.keys()].sort()).toEqual([...vi!.keys()].sort())
      for (const [key, text] of vi!) {
        expect(text.trim(), `vi:${ns}:${key}`).not.toBe('')
        expect(en!.get(key)?.trim(), `en:${ns}:${key}`).not.toBe('')
        expect(variables(en!.get(key)!), `${ns}:${key}`).toEqual(variables(text))
      }
    })
  }
})

// Văn bản được phép viết thẳng: tên sản phẩm.
const ALLOWED = new Set(['Crucible'])
const TEXT_ATTRIBUTES = new Set(['aria-label', 'title', 'placeholder', 'alt'])
const hasLetter = (s: string) => /\p{L}/u.test(s) && !ALLOWED.has(s.trim())

/** Tìm văn bản hiển thị viết thẳng trong JSX thay vì đi qua khóa dịch. */
export function findLiterals(fileName: string, source: string): string[] {
  const file = ts.createSourceFile(fileName, source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX)
  const found: string[] = []
  const report = (node: ts.Node, text: string) => {
    const { line } = file.getLineAndCharacterOfPosition(node.getStart())
    found.push(`${fileName}:${line + 1}: ${text.trim()}`)
  }
  const visit = (node: ts.Node) => {
    if (ts.isJsxText(node) && hasLetter(node.text)) report(node, node.text)
    if (ts.isJsxExpression(node) && node.expression && ts.isJsxElement(node.parent) && ts.isStringLiteralLike(node.expression) && hasLetter(node.expression.text)) {
      report(node, node.expression.text)
    }
    if (ts.isJsxAttribute(node) && TEXT_ATTRIBUTES.has(node.name.getText())) {
      const init = node.initializer
      const value = init && ts.isJsxExpression(init) ? init.expression : init
      if (value && ts.isStringLiteralLike(value) && hasLetter(value.text)) report(node, `${node.name.getText()}="${value.text}"`)
    }
    ts.forEachChild(node, visit)
  }
  visit(file)
  return found
}

describe('văn bản trong JSX đi qua khóa dịch (AC-P2.1, W3)', () => {
  it('phát hiện chữ viết thẳng, thuộc tính văn bản viết thẳng', () => {
    const sample = `export const A = () => (
      <div title="Tiêu đề" aria-label={'Đóng'} className="flex">
        Xin chào {'thế giới'} <span>Crucible</span> {t('x')} 42
      </div>
    )`
    expect(findLiterals('a.tsx', sample)).toEqual([
      'a.tsx:2: title="Tiêu đề"',
      `a.tsx:2: aria-label="Đóng"`,
      'a.tsx:3: Xin chào',
      'a.tsx:3: thế giới',
    ])
  })

  it('mã nguồn giao diện không có văn bản viết thẳng', () => {
    const files = import.meta.glob(['/src/**/*.tsx', '!/src/**/*.test.tsx'], { query: '?raw', import: 'default', eager: true })
    expect(Object.keys(files).length).toBeGreaterThan(0)
    const found = Object.entries(files).flatMap(([name, source]) => findLiterals(name, source as string))
    expect(found).toEqual([])
  })
})
