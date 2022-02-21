import React, {createContext, ReactChild, useEffect, useState} from 'react';
import './styles/App.css';
import {SocketProvider} from "./SocketProvider";
import GetName from "./components/GetName";
import {Socket} from "./components/Socket";

function App() {
    const [myName, setMyName] = useState();

    if (!myName) return (
        <GetName myName={myName} setMyName={setMyName}/>
    );

    return (
        <SocketProvider myName={myName}>
            <Socket/>
        </SocketProvider>
    );
}

export default App;
