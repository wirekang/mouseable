import mock from "./mock";

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
}

const isDev = process.env.NODE_ENV === "development";
if (isDev) {
  console.log("DEVELOPMENT MODE");
}

export async function loadBind(): Promise<GoBind> {
  if (isDev) {
    await new Promise((r) => setTimeout(r, 1000));
    return mock;
  }
  return window.__loadBind__();
}

export async function changeFunction(name: string, key: FunctionKey): Promise<boolean> {
  if (isDev) {
    mock.functionNameKeyMap[name] = key;
    await new Promise((r) => setTimeout(r, 1000));
    return true;
  }

  return window.__changeFunction__(name, key);
}
