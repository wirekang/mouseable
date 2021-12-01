import React from "react";
import { GoBind } from "./gobind";

const MyContext = React.createContext<{
  requestChangeFunctionKey: (name: string) => void;
}>({} as any);

export default MyContext;
