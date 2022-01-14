import React from "react";

interface Props {
  children: React.ReactNode;
}

export default function HalfRow(props: Props): JSX.Element {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-between",
        alignContent: "center",
        margin: 0,
        padding: 0,
      }}
    >
      {props.children}
    </div>
  );
}
