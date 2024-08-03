import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import UnoCSS from "unocss/vite";
import { presetUno, presetAttributify, presetIcons } from "unocss";
import Components from 'unplugin-vue-components/vite';
import {PrimeVueResolver} from '@primevue/auto-import-resolver';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Components({
      resolvers: [
        PrimeVueResolver()
      ]
    }),
    UnoCSS({
      presets: [presetUno(), presetAttributify(), presetIcons()],
    })
  ],
  server: {
    port: 9245, 
    strictPort: true, 
  }
});
function AutoImport(arg0: { resolvers: import("unplugin-vue-components/types").ComponentResolver[][]; dts: string; }): import("vite").PluginOption {
  throw new Error("Function not implemented.");
}

