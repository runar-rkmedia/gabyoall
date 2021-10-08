import { pnpPlugin } from '@yarnpkg/esbuild-plugin-pnp'
import { build } from 'esbuild'
import sveltePlugin from 'esbuild-svelte'
import sveltePreprocess from 'svelte-preprocess'
import fs from 'fs'
import path from 'path'
import { exec } from 'child_process'

const createTypescriptApiDefinitions = () => {
  exec('yarn gen', (err, out, errOut) => {
    if (err) {
      console.error('Failed to create typescript-defintitions for api: ', err)
    }
    if (errOut) {
      console.error(
        'Failed to create typescript-defintitions for api: ',
        errOut
      )
    }
    if (!err && !errOut) {
      console.log('Created typescript-defintions for api', out)
    }
  })
}
createTypescriptApiDefinitions()

const args = process.argv.slice(2)
const srcDir = './src/'
const outDir = './dist/'
const staticDir = './static/'

const isDev =
  args.includes('dev') || args.includes('-d') || args.includes('development')

if (!fs.existsSync(outDir)) {
  fs.mkdirSync(outDir)
}

await build({
  plugins: [
    pnpPlugin(),
    sveltePlugin({
      preprocess: sveltePreprocess(),
    }),
  ],
  entryPoints: [srcDir + 'entry.ts'],
  bundle: true,
  outdir: outDir,
  logLevel: 'info',
  sourcemap: 'external',
  legalComments: 'external',
  minify: true,
  ...(isDev && {
    watch: true,
    legalComments: 'none',
    minify: false,
    sourcemap: 'inline',
  }),
}).then(() => {
  if (isDev) {
    createTypescriptApiDefinitions()
  }
})

fs.copyFile(srcDir + 'index.html', outDir + '/index.html', (err) => {
  if (err) throw err
})
const staticFiles = fs.readdirSync('./static')
await Promise.all(
  staticFiles.map((f) => {
    return new Promise((res) => {
      const src = path.join(staticDir, f)
      const target = path.join(outDir, f)
      return fs.copyFile(src, target, res)
    })
  })
)
