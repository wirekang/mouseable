import React from "react";
import GroupBox from "./GroupBox";
import { openLink } from "../gobind";
import GitHubButton from "./GitHubButton";

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
      </div>
    </GroupBox>
  );
}
