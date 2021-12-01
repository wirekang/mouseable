import React from "react";
import { FunctionDefinition, FunctionNameKeyRecord } from "../gobind";
import Row from "./Row";
import GroupBox from "./GroupBox";
import FunctionRow from "./FunctionRow";

interface Props {
  def: FunctionDefinition[];
  record: FunctionNameKeyRecord;
}

export default function FunctionGroupBox(props: Props): JSX.Element {
  const categoryMap = new Map<string, FunctionDefinition[]>();
  props.def.forEach((fd) => {
    const arr = categoryMap.get(fd.Category);
    if (!arr) {
      categoryMap.set(fd.Category, [fd]);
    } else {
      arr.push(fd);
    }
  });
  return (
    <div>
      <GroupBox title="Function">
        <Row
          name="Name"
          column2={<span>Key</span>}
          column3={<span>When</span>}
          description="Description"
          style={{
            margin: "10px 0 5px 11px",
            fontWeight: "bold",
          }}
        />
        {Array.from(categoryMap, ([key, value]): [string, FunctionDefinition[]] => [key, value]).map(
          ([category, fds]) => (
            <GroupBox key={category} title={category}>
              {fds.map((fd) => (
                <FunctionRow key={fd.Name} def={fd} fKey={props.record[fd.Name]} />
              ))}
            </GroupBox>
          ),
        )}
      </GroupBox>
    </div>
  );
}
