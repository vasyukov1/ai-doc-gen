import { defineConfig } from "vite";

// Builds one stylesheet and one script with stable names, so Go templates can
// reference them directly and the binary can embed them.
export default defineConfig({
  build: {
    outDir: "static",
    emptyOutDir: true,
    assetsDir: "",
    rollupOptions: {
      input: "src/main.ts",
      output: {
        entryFileNames: "app.js",
        assetFileNames: (asset) => {
          const name = asset.names?.[0] ?? asset.name ?? "";
          return name.endsWith(".css") ? "app.css" : "assets/[name][extname]";
        },
      },
    },
  },
  server: {
    proxy: {
      "/": "http://localhost:8080",
    },
  },
});
