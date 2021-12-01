import React from "react";
import { DataDefinition, DataNameValueRecord } from "../gobind";
import GroupBox from "./GroupBox";
import DataRow from "./DataRow";
import Row from "./Row";

interface Props {
  def: DataDefinition[];
  record: DataNameValueRecord;
}

export default function DataGroupBox(props: Props): JSX.Element {
  return (
    <GroupBox title="Data">
      <Row
        name="Name"
        column2={<span>Value</span>}
        column3={<span>Type</span>}
        description="Description"
        style={{
          margin: "10px 0 5px 11px",
          fontWeight: "bold",
        }}
      />
      {props.def.map((dd) => (
        <DataRow key={dd.Name} def={dd} value={props.record[dd.Name]} />
      ))}
    </GroupBox>
  );
}
