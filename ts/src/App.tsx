import React, { useState } from "react";
import { changeFunction, DataDefinition, FunctionDefinition, FunctionKey, GoBind, loadBind } from "./gobind";
import { useAsync } from "react-use";
import FunctionGroupBox from "./components/FunctionGroupBox";
import FunctionKeyInput from "./components/FunctionKeyInput";
import Background from "./components/Background";
import MyContext from "./MyContext";
import DataGroupBox from "./components/DataGroupBox";
import InfoGroupBox from "./components/InfoGroupBox";

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

  const changeFunctionKeyInput = async (name: string, key: FunctionKey) => {
    const rst = await changeFunction(name, key);
    if (rst) {
      setKeyInputState({ isOpen: false });
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
        <InfoGroupBox version={goBind.version} />
        <FunctionGroupBox def={goBind.functionDefinitions} record={goBind.functionNameKeyMap} />
        <DataGroupBox def={goBind.dataDefinitions} record={goBind.dataNameValueMap} />

        {keyInputState.isOpen && keyInputState.name && keyInputState.key && (
          <Background onClick={closeFunctionKeyInput}>
            <FunctionKeyInput
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
