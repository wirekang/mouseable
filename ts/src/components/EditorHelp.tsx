import React from "react";

interface Props {}

export default function EditorHelp(props: Props): JSX.Element {
  return (
    <div
      style={{
        fontSize: 12,
        marginTop: 10,
        marginLeft: 10,
      }}
    >
      <ul>
        <li>Press Ctrl-I to show suggestions.</li>
        <li>Press F1 to insert key.</li>
        <li>Press F2 to save.</li>
      </ul>
    </div>
  );
}
