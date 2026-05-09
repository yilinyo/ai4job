import DOMPurify from 'dompurify'
// `marked` 4.x does not ship local type declarations in this project setup,
// so we intentionally consume its CommonJS entry to keep Vue 2/Webpack 4 happy.
// eslint-disable-next-line @typescript-eslint/no-var-requires
const { marked } = require('marked') as { marked: any }

const allowedTags = [
  'h1',
  'h2',
  'h3',
  'h4',
  'h5',
  'h6',
  'p',
  'br',
  'strong',
  'em',
  'blockquote',
  'ul',
  'ol',
  'li',
  'code',
  'pre',
  'hr',
  'table',
  'thead',
  'tbody',
  'tr',
  'th',
  'td',
  'a',
  'img'
]

const allowedAttributes = ['href', 'title', 'target', 'rel', 'src', 'alt']

marked.setOptions({
  gfm: true,
  breaks: true,
  mangle: false,
  headerIds: false
})

function escapeRawHtml(content: string) {
  return content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

export function renderResumeMarkdown(content: string) {
  if (!content.trim()) return ''

  const rawHtml = marked.parse(escapeRawHtml(content))

  return DOMPurify.sanitize(rawHtml, {
    ALLOWED_TAGS: allowedTags,
    ALLOWED_ATTR: allowedAttributes
  })
}
