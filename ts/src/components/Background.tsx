import React from "react";

interface Props {
  children: React.ReactNode;
  onClick?: () => void;
}

export default function Background(props: Props): JSX.Element {
  return (
    <div
      onClick={props.onClick}
      style={{
        position: "absolute",
        left: 0,
        right: 0,
        top: 0,
        bottom: 0,
        backgroundColor: "#000000aa",
      }}
    >
      <div>{props.children}</div>
    </div>
  );
}
