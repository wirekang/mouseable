import React from "react";
import { FunctionKey } from "../gobind";
import MyContext from "../MyContext";
import { functionKeyToText } from "../util/function";
import { useAsync } from "react-use";

interface Props {
  name: string;
  functionKey: FunctionKey;
}

export default function KeyHolder(props: Props): JSX.Element {
  const text = useAsync(functionKeyToText.bind(null, props.functionKey));
  return (
    <MyContext.Consumer>
      {(v) => (
        <div
          style={{
            cursor: "pointer",
            backgroundColor: "#f8f8f8",
            border: "1px solid #ccc",
            fontSize: 10,
            margin: "0 2px",
            padding: "0 2px",
            height: "14px",
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
          }}
          onClick={() => {
            v.requestChangeFunctionKey(props.name);
          }}
        >
          {text.value}
        </div>
      )}
    </MyContext.Consumer>
  );
}
