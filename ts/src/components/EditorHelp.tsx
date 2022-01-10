import React from "react";

interface Props {}

export default function EditorHelp(props: Props): JSX.Element {
  return (
    <div
      style={{
        fontSize: 11,
        marginTop: 10,
        marginLeft: 10,
      }}
    >
      <p>Press F1 to insert key.</p>
    </div>
  );
}
