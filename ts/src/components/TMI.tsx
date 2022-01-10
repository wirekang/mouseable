import React, { useState } from "react";
import TextButton from "./TextButton";

interface Props {}

export default function TMI(props: Props): JSX.Element {
  const [isOpen, setIsOpen] = useState(false);
  return (
    <div>
      <TextButton
        style={{ fontSize: 12 }}
        onClick={() => {
          setIsOpen((b) => !b);
        }}
      >
        TooMuchInformation
      </TextButton>
      {isOpen && <p>Thank you!!</p>}
    </div>
  );
}
