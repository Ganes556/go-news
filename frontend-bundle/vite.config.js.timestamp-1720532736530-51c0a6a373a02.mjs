// vite.config.js
import inject from "file:///D:/golang/news-golang/frontend-bundle/node_modules/@rollup/plugin-inject/dist/es/index.js";
import { createRequire } from "node:module";
import ckeditor5 from "file:///D:/golang/news-golang/frontend-bundle/node_modules/@ckeditor/vite-plugin-ckeditor5/dist/index.mjs";
import { defineConfig } from "file:///D:/golang/news-golang/frontend-bundle/node_modules/vite/dist/node/index.js";
import postcssNesting from "file:///D:/golang/news-golang/frontend-bundle/node_modules/postcss-nesting/dist/index.mjs";
var __vite_injected_original_import_meta_url = "file:///D:/golang/news-golang/frontend-bundle/vite.config.js";
var require2 = createRequire(__vite_injected_original_import_meta_url);
var vite_config_default = defineConfig({
  plugins: [
    inject({
      htmx: "htmx.org"
      // htmx does use eval, based on this issue: https://github.com/bigskysoftware/htmx/issues/1944
    }),
    ckeditor5({ theme: require2.resolve("@ckeditor/ckeditor5-theme-lark") })
  ],
  css: {
    postcss: {
      plugins: [postcssNesting]
      // use this, to be able to use nested css
    }
  },
  server: {
    origin: "http://localhost:5173"
  },
  build: {
    manifest: true,
    rollupOptions: {
      input: "main.js",
      output: {
        format: "iife",
        dir: "../static",
        entryFileNames: "main.js"
      }
    }
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJEOlxcXFxnb2xhbmdcXFxcbmV3cy1nb2xhbmdcXFxcZnJvbnRlbmQtYnVuZGxlXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxnb2xhbmdcXFxcbmV3cy1nb2xhbmdcXFxcZnJvbnRlbmQtYnVuZGxlXFxcXHZpdGUuY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi9nb2xhbmcvbmV3cy1nb2xhbmcvZnJvbnRlbmQtYnVuZGxlL3ZpdGUuY29uZmlnLmpzXCI7aW1wb3J0IGluamVjdCBmcm9tICdAcm9sbHVwL3BsdWdpbi1pbmplY3QnO1xyXG5pbXBvcnQgeyBjcmVhdGVSZXF1aXJlIH0gZnJvbSAnbm9kZTptb2R1bGUnO1xyXG5pbXBvcnQgY2tlZGl0b3I1IGZyb20gJ0Bja2VkaXRvci92aXRlLXBsdWdpbi1ja2VkaXRvcjUnO1xyXG5pbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tICd2aXRlJztcclxuaW1wb3J0IHBvc3Rjc3NOZXN0aW5nIGZyb20gJ3Bvc3Rjc3MtbmVzdGluZyc7XHJcbmNvbnN0IHJlcXVpcmUgPSBjcmVhdGVSZXF1aXJlKGltcG9ydC5tZXRhLnVybCk7XHJcblxyXG4vKiogQHR5cGUge2ltcG9ydCgndml0ZScpLlVzZXJDb25maWd9ICovXHJcbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XHJcbiAgcGx1Z2luczogW1xyXG4gICAgaW5qZWN0KHtcclxuICAgICAgaHRteDogJ2h0bXgub3JnJywgLy8gaHRteCBkb2VzIHVzZSBldmFsLCBiYXNlZCBvbiB0aGlzIGlzc3VlOiBodHRwczovL2dpdGh1Yi5jb20vYmlnc2t5c29mdHdhcmUvaHRteC9pc3N1ZXMvMTk0NFxyXG4gICAgfSksXHJcbiAgICBja2VkaXRvcjUoeyB0aGVtZTogcmVxdWlyZS5yZXNvbHZlKCdAY2tlZGl0b3IvY2tlZGl0b3I1LXRoZW1lLWxhcmsnKSB9KSxcclxuICBdLFxyXG4gIGNzczoge1xyXG4gICAgcG9zdGNzczoge1xyXG4gICAgICBwbHVnaW5zOiBbcG9zdGNzc05lc3RpbmddLCAvLyB1c2UgdGhpcywgdG8gYmUgYWJsZSB0byB1c2UgbmVzdGVkIGNzc1xyXG4gICAgfSxcclxuICB9LFxyXG4gIHNlcnZlcjoge1xyXG4gICAgb3JpZ2luOiAnaHR0cDovL2xvY2FsaG9zdDo1MTczJyxcclxuICB9LFxyXG4gIGJ1aWxkOiB7XHJcbiAgICBtYW5pZmVzdDogdHJ1ZSxcclxuICAgIHJvbGx1cE9wdGlvbnM6IHtcclxuICAgICAgaW5wdXQ6ICdtYWluLmpzJyxcclxuICAgICAgb3V0cHV0OiB7XHJcbiAgICAgICAgZm9ybWF0OiAnaWlmZScsXHJcbiAgICAgICAgZGlyOiAnLi4vc3RhdGljJyxcclxuICAgICAgICBlbnRyeUZpbGVOYW1lczogJ21haW4uanMnLFxyXG4gICAgICB9LFxyXG4gICAgfSxcclxuICB9LFxyXG59KTtcclxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUF5UyxPQUFPLFlBQVk7QUFDNVQsU0FBUyxxQkFBcUI7QUFDOUIsT0FBTyxlQUFlO0FBQ3RCLFNBQVMsb0JBQW9CO0FBQzdCLE9BQU8sb0JBQW9CO0FBSjhKLElBQU0sMkNBQTJDO0FBSzFPLElBQU1BLFdBQVUsY0FBYyx3Q0FBZTtBQUc3QyxJQUFPLHNCQUFRLGFBQWE7QUFBQSxFQUMxQixTQUFTO0FBQUEsSUFDUCxPQUFPO0FBQUEsTUFDTCxNQUFNO0FBQUE7QUFBQSxJQUNSLENBQUM7QUFBQSxJQUNELFVBQVUsRUFBRSxPQUFPQSxTQUFRLFFBQVEsZ0NBQWdDLEVBQUUsQ0FBQztBQUFBLEVBQ3hFO0FBQUEsRUFDQSxLQUFLO0FBQUEsSUFDSCxTQUFTO0FBQUEsTUFDUCxTQUFTLENBQUMsY0FBYztBQUFBO0FBQUEsSUFDMUI7QUFBQSxFQUNGO0FBQUEsRUFDQSxRQUFRO0FBQUEsSUFDTixRQUFRO0FBQUEsRUFDVjtBQUFBLEVBQ0EsT0FBTztBQUFBLElBQ0wsVUFBVTtBQUFBLElBQ1YsZUFBZTtBQUFBLE1BQ2IsT0FBTztBQUFBLE1BQ1AsUUFBUTtBQUFBLFFBQ04sUUFBUTtBQUFBLFFBQ1IsS0FBSztBQUFBLFFBQ0wsZ0JBQWdCO0FBQUEsTUFDbEI7QUFBQSxJQUNGO0FBQUEsRUFDRjtBQUNGLENBQUM7IiwKICAibmFtZXMiOiBbInJlcXVpcmUiXQp9Cg==
