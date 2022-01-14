import React from "react";

interface Props {}

export default function TerminateButton(props: Props): JSX.Element {
  const onClick = () => {
    window.__terminate();
    window.close();
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
