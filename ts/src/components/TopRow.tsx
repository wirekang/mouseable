import React from "react";
import GitHubButton from "./GitHubButton";
import TerminateButton from "./TerminateButton";

interface Props {
  version: string;
}

export default function TopRow(props: Props): JSX.Element {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        margin: 3,
      }}
    >
      <span>Version: {props.version}</span>
      <GitHubButton />
      <TerminateButton />
    </div>
  );
}
