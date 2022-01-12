declare global {
  interface Window {
    __ping__: () => Promise<number>;
    __getVersion__: () => Promise<string>;
    __getSchema__: () => Promise<string>;
    __getNextKey__: () => Promise<string>;
    __terminate__: () => Promise<void>;
    __openLink__: (url: string) => Promise<void>;
    __getConfigNames__: () => Promise<string[]>;
    __getConfig__: (name: string) => Promise<string>;
    __saveConfig__: (json: string) => Promise<void>;
  }
}
export const isDev = process.env.NODE_ENV === "development";
if (isDev) {
  window.__ping__ = () => resolve(1);
  window.__getSchema__ = () => resolve("{}");
  window.__getNextKey__ = () => resolve("Shift-Q");
  window.__terminate__ = () => resolve(undefined);
  window.__getVersion__ = () => resolve("dd.ee.vv");
  window.__openLink__ = () => resolve(undefined);
  window.__getConfigNames__ = () => resolve(["config-1.json", "my-config.json", "your-config.json"]);
  window.__getConfig__ = (name) => resolve(`{"name": "${name}"}`);
  window.__saveConfig__ = () => resolve(undefined);
}

function resolve<T>(value: T, ms = 500): Promise<T> {
  return new Promise((r) => {
    setTimeout(() => {
      r(value);
    }, ms);
  });
}
