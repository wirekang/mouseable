import React from "react";

interface Props {
  name: string;
  column2: React.ReactNode;
  column3?: React.ReactNode;
  description: string;
  style?: React.CSSProperties;
}

export default function Row(props: Props): JSX.Element {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        width: "100%",
        height: 16,
        margin: "1px 5px",
        fontSize: 11,
        ...props.style,
      }}
    >
      <span
        style={{
          fontWeight: "bold",
          minWidth: 130,
          maxWidth: 130,
        }}
      >
        {props.name}
      </span>

      <div
        style={{
          maxWidth: 150,
          minWidth: 150,
        }}
      >
        {props.column2}
      </div>

      {props.column3 && (
        <div
          style={{
            maxWidth: 70,
            minWidth: 70,
          }}
        >
          {props.column3}
        </div>
      )}
      <span
        style={{
          whiteSpace: "nowrap",
          width: "min-content",
          overflowX: "hidden",
          flexGrow: 1,
        }}
      >
        {props.description}
      </span>
    </div>
  );
}
