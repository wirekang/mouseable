import React, { useRef, useState } from "react";
import { useClickAway } from "react-use";
import { FunctionDefinition, FunctionKey } from "../gobind";
import Checkbox from "./Checkbox";
import { functionKeyToText } from "../util/function";
import InputKeyCode from "./InputKeyCode";

interface Props {
  name?: string;
  fKey?: FunctionKey;
  close: () => void;
  change: (name: string, fKey: FunctionKey) => void;
}

export default function FunctionKeyInput(props: Props): JSX.Element {
  const [key, setKey] = useState(props.fKey!);
  if (!props.name || !key) {
    return <>ERROR</>;
  }

  return (
    <div
      style={{
        position: "absolute",
        left: 0,
        right: 0,
        top: 0,
        bottom: 0,
        margin: "auto",
        width: 300,
        height: 140,
        backgroundColor: "white",
        border: " 1px solid black",
        zIndex: 100,
        display: "flex",
        flexDirection: "column",
        justifyContent: "flex-start",
        alignItems: "center",
      }}
      onClick={(e) => {
        e.stopPropagation();
      }}
    >
      <span style={{ fontSize: 24, fontWeight: "bold", marginTop: 6 }}>{props.name}</span>
      <span style={{ fontSize: 12, margin: "15px 0" }}>{functionKeyToText(key)}</span>
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
          width: "95%",
          alignItems: "center",
        }}
      >
        {(["IsWin", "IsControl", "IsAlt", "IsShift"] as const).map((k) => (
          <Checkbox
            style={{
              fontSize: 16,
            }}
            isChecked={key[k]}
            onChange={(v) => {
              setKey((p) => ({ ...p, [k]: v }));
            }}
          >
            {k.replace("Is", "")}
          </Checkbox>
        ))}
      </div>
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
          width: "95%",
          alignItems: "center",
          margin: "10px 0",
        }}
      >
        <InputKeyCode keyCode={key.KeyCode} onChange={() => {}} />
        <Checkbox
          style={{
            fontSize: 13,
          }}
          isChecked={key.IsDouble}
          onChange={(IsDouble) => {
            setKey((v) => ({ ...v, IsDouble }));
          }}
        >
          Double Press
        </Checkbox>
      </div>
    </div>
  );
}
