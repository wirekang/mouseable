import React, { useState } from "react";
import { DataDefinition, FunctionDefinition, FunctionKey, GoBind, loadBind } from "./gobind";
import { useAsync } from "react-use";
import FunctionCategories from "./components/FunctionCategories";
import FunctionKeyInput from "./components/FunctionKeyInput";
import Background from "./components/Background";
import MyContext from "./MyContext";

function App() {
  const [requesterCounter, setRequesterCounter] = useState(0);
  const goBindState = useAsync(loadBind, [requesterCounter]);
  const [keyInputState, setKeyInputState] = useState<{
    isOpen: boolean;
    name?: string;
    key?: FunctionKey;
  }>({
    isOpen: false,
  });

  if (goBindState.loading || !goBindState.value) {
    return <h4>Loading...</h4>;
  }

  const goBind = goBindState.value!;

  const requestChangeFunctionKey = (name: string) => {
    const key = goBind.functionNameKeyMap[name];
    setKeyInputState({ name, isOpen: true, key });
  };

  const closeFunctionKeyInput = () => {
    setKeyInputState({ isOpen: false, key: undefined, name: undefined });
  };

  const changeFunctionKeyInput = async (name: string, fKey: FunctionKey) => {
    const rst = await window.__changeFunction__(name, fKey);
    if (rst) {
      setRequesterCounter((v) => v + 1);
    }
  };

  return (
    <MyContext.Provider
      value={{
        requestChangeFunctionKey,
      }}
    >
      <div>
        <FunctionCategories def={goBind.functionDefinitions} record={goBind.functionNameKeyMap} />

        {keyInputState.isOpen && (
          <Background close={closeFunctionKeyInput}>
            <FunctionKeyInput
              close={closeFunctionKeyInput}
              name={keyInputState.name}
              fKey={keyInputState.key}
              change={changeFunctionKeyInput}
            />
          </Background>
        )}
      </div>
    </MyContext.Provider>
  );
}

export default App;
