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

export async function loadBind(): Promise<GoBind> {
  const f = (window as any).__loadBind;
  if (!f) {
    console.log("USE MOCK");
    return Promise.resolve(mock);
  }
  return f();
}
