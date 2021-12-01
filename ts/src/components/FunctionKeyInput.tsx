import React, { useState } from "react";
import { FunctionKey } from "../gobind";
import Checkbox from "./Checkbox";
import { functionKeyToText } from "../util/function";
import InputKeyCode from "./InputKeyCode";

interface Props {
  name: string;
  fKey: FunctionKey;
  change: (name: string, fKey: FunctionKey) => void;
}

export default function FunctionKeyInput(props: Props): JSX.Element {
  const [bufferKey, setBufferKey] = useState(props.fKey);

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
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
          alignItems: "flex-end",
          width: "95%",
        }}
      >
        <button
          onClick={() => {
            setBufferKey({
              KeyCode: 0,
              IsDouble: false,
              IsWin: false,
              IsShift: false,
              IsAlt: false,
              IsControl: false,
            });
          }}
        >
          Unset
        </button>
        <span style={{ fontSize: 18, fontWeight: "bold", marginTop: 6 }}>{props.name}</span>
        <button
          onClick={() => {
            props.change(props.name, bufferKey);
          }}
        >
          Apply
        </button>
      </div>
      <span style={{ fontSize: 12, margin: "15px 0" }}>{functionKeyToText(bufferKey)}</span>
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
            key={k}
            style={{
              fontSize: 16,
            }}
            isChecked={bufferKey[k]}
            onChange={(v) => {
              setBufferKey((p) => ({ ...p, [k]: v }));
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
        <InputKeyCode
          keyCode={bufferKey.KeyCode}
          onChange={(KeyCode) => {
            setBufferKey((p) => ({ ...p, KeyCode }));
          }}
        />
        <Checkbox
          style={{
            fontSize: 13,
          }}
          isChecked={bufferKey.IsDouble}
          onChange={(IsDouble) => {
            setBufferKey((p) => ({ ...p, IsDouble }));
          }}
        >
          Double Press
        </Checkbox>
      </div>
    </div>
  );
}
