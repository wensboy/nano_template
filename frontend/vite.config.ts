import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath } from "node:url";

const apiProxyTarget = "http://localhost:3000";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: apiProxyTarget,
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, "/api"), // 重写 uri
        configure: (proxy, options) => {
          proxy.on("proxyReq", (proxyReq, req) => {
						const proxiedPath =
								typeof proxyReq.path === "string" ? proxyReq.path : req.url ?? "";
            console.log(
              "[proxy] - ",
              req.method,
              req.url,
              " -> ",
              `${apiProxyTarget}${proxiedPath}`,
            );
          });
        },
      },
    },
  },
});
