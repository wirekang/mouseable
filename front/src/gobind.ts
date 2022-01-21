declare global {
  interface Window {
    __ping: () => Promise<number>;
    __getVersion: () => Promise<string>;
    __loadSchema: () => Promise<string>;
    __getNextKey: () => Promise<string>;
    __terminate: () => Promise<void>;
    __openLink: (url: string) => Promise<void>;
    __loadConfigNames: () => Promise<string[]>;
    __loadAppliedConfigName: () => Promise<string>;
    __loadConfig: (name: string) => Promise<string>;
    __saveConfig: (name: string, json: string) => Promise<void>;
    __applyConfig: (name: string) => Promise<void>;
  }
}

export const isDev = process.env.NODE_ENV === "development";

if (isDev) {
  const configs = ["config-1.json", "cfg2.json", "long-long-config-name.json"];
  const configJsons = new Map<string, string>();
  let appliedConfigName = "cfg2.json";

  window.__ping = () => resolve(1);
  window.__loadSchema = () => resolve("{}");
  window.__getNextKey = () => resolve("Shift+Q");
  window.__terminate = () => resolve(undefined);
  window.__getVersion = () => resolve("dd.ee.vv");
  window.__openLink = () => resolve(undefined);
  window.__loadConfigNames = () => resolve(configs);
  window.__loadAppliedConfigName = () => resolve(appliedConfigName);
  window.__loadConfig = (name) => resolve(configJsons.get(name) ?? `${name} default`);
  window.__saveConfig = (name, json) =>
    resolve(() => {
      configs.push(name);
      configJsons.set(name, json);
    });
  window.__applyConfig = (name) =>
    resolve(() => {
      appliedConfigName = name;
    });
}

function resolve<T>(value: T | (() => T)): Promise<T> {
  return new Promise((r) => {
    setTimeout(() => {
      if (typeof value === "function") {
        r((value as () => T)());
      } else {
        r(value);
      }
    }, Math.random() * 1000);
  });
}
