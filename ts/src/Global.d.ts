import { FunctionKey, GoBind } from "./gobind";

declare global {
  interface Window {
    __loadBind__: () => Promise<GoBind>;
    __changeFunction__: (name: string, key: FunctionKey) => Promise<boolean>;
  }
}
