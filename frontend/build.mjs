import { pnpPlugin } from '@yarnpkg/esbuild-plugin-pnp'
import { build } from 'esbuild'
import sveltePlugin from 'esbuild-svelte'
import sveltePreprocess from 'svelte-preprocess'
import fs from 'fs'
import path from 'path'
import { exec } from 'child_process'
import { stdout } from 'process'

const args = process.argv.slice(2)
const srcDir = './src/'
const outDir = './dist/'
const staticDir = './static/'

const execP = (args) =>
  new Promise((res) =>
    exec(args, (err, out, stdErr) =>
      res([out, (err || stdErr) && { err, stdErr }])
    )
  )

const isDev =
  args.includes('dev') || args.includes('-d') || args.includes('development')
const withDTS =
  args.includes('dts') || args.includes('-t') || args.includes('types')

if (!fs.existsSync(outDir)) {
  fs.mkdirSync(outDir)
}

const createTypescriptApiDefinitions = async () => {
  if (!withDTS) {
    return
  }
  console.time('creating typescript api defintions...')
  const [res, err] = await execP('yarn gen')
  if (err) {
    console.error('Failed to create typescript-defintitions for api: ', err)
  } else {
    console.log('Created typescript-defintions for api', out)
  }
  console.timeEnd('creating typescript api defintions...')
}
const typecheck = async () => {
  console.time('typechecking')
  const [res, err] = await execP('yarn tsc --noEmit')
  if (res) {
    console.log(res)
  }
  if (err && !(res || '').includes('error')) {
    console.error(err)
  }
  console.timeEnd('typechecking')
}
createTypescriptApiDefinitions()
typecheck()

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
    watch: {
      onRebuild: (error, result) => {
        if (error) {
          console.error('watch build failed:', error)
        } else {
          console.log('watch build succeeded:', result)
        }
        createTypescriptApiDefinitions()
        typecheck()
      },
    },
    legalComments: 'none',
    minify: false,
    sourcemap: 'inline',
  }),
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
