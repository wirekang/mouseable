import React from "react";

interface Props {
  onClick: () => void
  children: React.ReactText
  style?: React.CSSProperties
}

export default function TextButton(props: Props): JSX.Element {
  return (
      <span style={{
        cursor: "pointer",
        textShadow: "1px 1px 2px grey",
        ...props.style
      }} onClick={props.onClick}>{props.children}</span>
  );
}