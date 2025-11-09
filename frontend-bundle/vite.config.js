import inject from '@rollup/plugin-inject';
import { createRequire } from 'node:module';
import ckeditor5 from '@ckeditor/vite-plugin-ckeditor5';
import { defineConfig } from 'vite';
import postcssNesting from 'postcss-nesting';
const require = createRequire(import.meta.url);

/** @type {import('vite').UserConfig} */
export default defineConfig({
  base: "/static/",
  plugins: [
    inject({
      htmx: 'htmx.org', // htmx does use eval, based on this issue: https://github.com/bigskysoftware/htmx/issues/1944
    }),
    ckeditor5({ theme: require.resolve('@ckeditor/ckeditor5-theme-lark') }),
  ],
  css: {
    postcss: {
      plugins: [postcssNesting], // use this, to be able to use nested css
    },
  },
  server: {
    origin: 'http://localhost:5173',
  },
  build: {
    manifest: true,
    rollupOptions: {
      input: 'main.js',
      output: {
        format: 'iife',
        dir: '../static',
        entryFileNames: 'main.js',
      },
    },
  },
});
