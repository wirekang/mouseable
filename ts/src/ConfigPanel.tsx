import React, { useEffect, useRef } from "react";

interface Props {
  appliedConfigName?: string;
  configNames?: string[];
  loadedConfigName?: string;
  onLoadConfig: (name: string) => void;
}

export default function ConfigPanel(props: Props): JSX.Element {
  return (
    <div style={{}}>
      <select
        value={props.loadedConfigName}
        onChange={(e) => {
          props.onLoadConfig(e.target.value);
        }}
      >
        {props.configNames?.map((name) => {
          return <option key={name}>{name}</option>;
        })}
      </select>
      <br />
      Applied: {props.appliedConfigName}
    </div>
  );
}
