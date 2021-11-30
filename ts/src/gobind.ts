import mock from "./mock";

interface FunctionDefinition {
  Name: string
  Category: string
  Description: string
  When: When
  Order: number
}

interface FunctionKey {
  IsDouble: boolean
  IsAlt: boolean
  IsControl: boolean
  IsShift: boolean
  IsWin: boolean
  KeyCode: number
}

interface DataDefinition {
  Name: string
  Description: string
  Type: DataType
}

enum DataType {
  Int = 0,
  Float = 1,
  Bool = 2,
  String = 3,
}

enum When {
  Activated = 0,
  Deactivated = 1,
  Any = 2,
}

export interface GoBind {
  functionDefinitions: FunctionDefinition[]
  dataDefinitions: DataDefinition[]
  functionNameKeyMap: Record<string, FunctionKey>
  dataNameValueMap: Record<string, string>
}

export async function loadBind(): Promise<GoBind> {
  const f = (window as any).__loadBind
  if (!f) {
    console.log("USE MOCK")
    return Promise.resolve(mock)
  }
  return f()
};