declare global {
  interface Window {
    __ping__: () => Promise<number>;
    __getVersion__: () => Promise<string>;
    __getSchema__: () => Promise<string>;
    __terminate__: () => Promise<void>;
    __openLink__: (url: string) => Promise<void>;
    __getConfigNames__: () => Promise<string[]>;
    __getConfig__: (name: string) => Promise<string>;
    __saveConfig__: (json: string) => Promise<void>;
  }
}
export const isDev = process.env.NODE_ENV === "development";
if (isDev) {
  window.__ping__ = () => Promise.resolve(1);
  window.__getSchema__ = () => Promise.resolve("{}");
  window.__terminate__ = () => Promise.resolve();
  window.__getVersion__ = () => Promise.resolve("dd.ee.vv");
  window.__openLink__ = () => Promise.resolve();
  window.__getConfigNames__ = () => Promise.resolve(["config-1.json", "my-config.json", "your-config.json"]);
  window.__getConfig__ = (name) => Promise.resolve(`{"name": "${name}"}`);
  window.__saveConfig__ = () => Promise.resolve();
}
