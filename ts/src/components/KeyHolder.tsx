import React, { useRef, useState } from "react";
import { FunctionKey } from "../gobind";
import { useClickAway, useMouse } from "react-use";
import ReactDOM from "react-dom";
import FunctionKeyInput from "./FunctionKeyInput";
import MyContext from "../MyContext";
import { functionKeyToText } from "../util/function";

interface Props {
  name: string;
  functionKey: FunctionKey;
}

export default function KeyHolder(props: Props): JSX.Element {
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
          {functionKeyToText(props.functionKey)}
        </div>
      )}
    </MyContext.Consumer>
  );
}
