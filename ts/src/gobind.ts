import mock from "./mock";

declare global {
  interface Window {
    __loadBind__: () => Promise<GoBind>;
    __getKeyCode__: () => Promise<number>;
    __changeFunction__: (name: string, key: FunctionKey) => Promise<boolean>;
    __changeData__: (name: string, value: string) => Promise<boolean>;
    __openLink__: (url: string) => Promise<void>;
  }
}

export interface FunctionDefinition {
  Name: string;
  Category: string;
  Description: string;
  When: When;
  Order: number;
}

export interface FunctionKey {
  IsDouble: boolean;
  IsAlt: boolean;
  IsControl: boolean;
  IsShift: boolean;
  IsWin: boolean;
  KeyCode: number;
}

export interface DataDefinition {
  Name: string;
  Description: string;
  Type: DataType;
}

export enum DataType {
  Int = 0,
  Float = 1,
  Bool = 2,
  String = 3,
}

export enum When {
  Activated = 0,
  Deactivated = 1,
  Any = 2,
}

export type FunctionNameKeyRecord = Record<string, FunctionKey>;
export type DataNameValueRecord = Record<string, string>;

export interface GoBind {
  functionDefinitions: FunctionDefinition[];
  dataDefinitions: DataDefinition[];
  functionNameKeyMap: FunctionNameKeyRecord;
  dataNameValueMap: DataNameValueRecord;
  version: string;
}

const isDev = process.env.NODE_ENV === "development";
if (isDev) {
  console.log("DEVELOPMENT MODE");
}

export async function loadBind(): Promise<GoBind> {
  if (isDev) {
    await new Promise((r) => setTimeout(r, 100));
    return mock;
  }
  return window.__loadBind__();
}

export async function changeFunction(name: string, key: FunctionKey): Promise<boolean> {
  if (isDev) {
    mock.functionNameKeyMap[name] = key;
    await new Promise((r) => setTimeout(r, 100));
    return true;
  }

  return window.__changeFunction__(name, key);
}

export async function changeData(name: string, value: string): Promise<boolean> {
  if (isDev) {
    mock.dataNameValueMap[name] = value;
    await new Promise((r) => setTimeout(r, 100));
    return true;
  }

  return window.__changeData__(name, value);
}

export async function getKeyCode(): Promise<number> {
  if (isDev) {
    await new Promise((r) => setTimeout(r, 100));
    const v = Math.round(Math.random() * 40 + 48);
    console.log("random keycode ", v.toString(16));
    return v;
  }

  return window.__getKeyCode__();
}

export async function openLink(url: string): Promise<void> {
  if (isDev) {
    return;
  }

  return window.__openLink__(url);
}
