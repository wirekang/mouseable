import React from "react";

const MyContext = React.createContext<{
  requestChangeFunctionKey: (name: string) => void;
}>({} as any);

export default MyContext;
