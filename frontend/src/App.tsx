import React, {createContext, ReactChild, useEffect, useState} from 'react';
import './styles/App.css';
// import Main from "./components/Main";
import GetName from "./components/GetName";
import {SocketProvider} from "./SocketProvider";
import {Socket} from "./components/Socket";

function App() {
    const [myName, setMyName] = useState();

    if (!myName) return (
        <GetName myName={myName} setMyName={setMyName}/>
    );

    return (
        <SocketProvider>
            <Socket myName={myName}/>
        </SocketProvider>
    );
}

export default App;
