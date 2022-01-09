import React from "react";
import GroupBox from "./GroupBox";
import GitHubButton from "./GitHubButton";
import TerminateButton from "./TerminateButton";

interface Props {
  version: string;
}

export default function InfoGroupBox(props: Props): JSX.Element {
  return (
    <GroupBox title="Mouseable">
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
          justifyContent: "space-between",
          margin: 3,
        }}
      >
        <span>{props.version}</span>
        <GitHubButton />
        <TerminateButton />
      </div>
    </GroupBox>
  );
}
