// vite.config.js
export default {
  build: {
    lib: {
      entry: 'cmd/static/js/app.js',
      name: 'NViro',
      fileName: (format) => `nviro.${format}.js`
    },
    outDir: 'cmd/static/dist/',
  }
}
