import React, {useState} from 'react';
import './styles/App.css';
import Main from "./components/Main";
import GetName from "./components/GetName";

function App() {
    const [myName, setMyName] = useState();
    console.log(myName);
    return (
        <div className="App">
            {
                myName ? <Main/> : <GetName myName={myName} setMyName={setMyName}/>
            }
        </div>
    );
}

export default App;
