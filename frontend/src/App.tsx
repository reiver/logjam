import React, {useEffect, useState} from 'react';
import './styles/App.css';
import Main from "./components/Main";
import GetName from "./components/GetName";
import useSocket from "./hooks/useSocket";

function App() {
    const [myName, setMyName] = useState();

    useEffect(()=>{
        console.log('myName is set to:', myName);
    }, myName);

    return (
        <div>
            {
                myName ?
                    <Main myName={myName}/> :
                    <GetName myName={myName} setMyName={setMyName}/>
            }
        </div>
    );
}

export default App;
