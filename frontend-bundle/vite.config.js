import inject from '@rollup/plugin-inject';
import { createRequire } from 'node:module';
const require = createRequire(import.meta.url);
import ckeditor5 from '@ckeditor/vite-plugin-ckeditor5';

/** @type {import('vite').UserConfig} */
export default {
  plugins: [
    inject({
      htmx: 'htmx.org',
    }),
    ckeditor5({ theme: require.resolve('@ckeditor/ckeditor5-theme-lark') }),
  ],
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
};
