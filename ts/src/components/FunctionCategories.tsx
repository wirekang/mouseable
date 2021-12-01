import React from "react";
import { FunctionDefinition, FunctionNameKeyRecord, When } from "../gobind";
import Row from "./Row";
import GroupBox from "./GroupBox";
import KeyHolder from "./KeyHolder";

interface Props {
  def: FunctionDefinition[];
  record: FunctionNameKeyRecord;
}

export default function FunctionCategories(props: Props): JSX.Element {
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
                <Row
                  key={fd.Name}
                  name={fd.Name}
                  description={fd.Description}
                  column2={<KeyHolder name={fd.Name} functionKey={props.record[fd.Name]} />}
                  column3={
                    <span
                      style={{
                        fontSize: 8,
                        transform: "scaleX(0.6)",
                      }}
                    >
                      {When[fd.When]}
                    </span>
                  }
                />
              ))}
            </GroupBox>
          ),
        )}
      </GroupBox>
    </div>
  );
}
