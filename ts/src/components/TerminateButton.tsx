import React from "react";
import { terminate } from "../gobind";

interface Props {}

export default function TerminateButton(props: Props): JSX.Element {
  const onClick = () => {
    terminate();
  };
  return (
    <div style={{ fontSize: 10 }}>
      You can close this window safely.
      <button style={{ fontSize: 12 }} onClick={onClick}>
        Terminate Program
      </button>
    </div>
  );
}
