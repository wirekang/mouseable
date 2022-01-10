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
const isDev = process.env.NODE_ENV === "development";
if (isDev) {
  console.log("DEVELOPMENT MODE");
}

export async function ping(): Promise<boolean> {
  console.log("ping");
  if (isDev) {
    await new Promise((r) => {
      setTimeout(r, 100);
    });
    return true;
  }

  const r = await window.__ping__();
  console.log(r);
  return !!r;
}

export function terminate() {
  console.log("terminate");
  if (isDev) {
    window.close();
    return;
  }

  window.__terminate__().then(window.close);
  return;
}

export async function getVersion(): Promise<string> {
  console.log("getVersion");
  if (isDev) {
    return "d.e.v";
  }

  const r = await window.__getVersion__();
  console.log(r);
  return r;
}

export async function openLink(url: string): Promise<void> {
  console.log(`openLink ${url}`);
  if (isDev) {
    return;
  }

  const r = await window.__openLink__(url);
  console.log(r);
  return r;
}

export async function getConfigNames(): Promise<string[]> {
  console.log("getConfigNames");
  if (isDev) {
    return ["config-1.json", "my-config.json", "your-config.json"];
  }

  const r = await window.__getConfigNames__();
  console.log(r);
  return r;
}

export async function getConfig(name: string): Promise<string> {
  console.log(`getConfig ${name}`);
  if (isDev) {
    return `{"text":"thsi is test"}`;
  }

  const r = await window.__getConfig__(name);
  console.log(r);
  return r;
}

export async function getSchema(): Promise<string> {
  console.log("getSchema");
  if (isDev) {
    return "{}";
  }

  const r = await window.__getSchema__();
  console.log(r);
  return r;
}

export async function saveConfig(json: string): Promise<void> {
  console.log(`saveConfig ${json}`);
  if (isDev) {
    return;
  }

  const r = await window.__saveConfig__(json);
  console.log(r);
  return r;
}
