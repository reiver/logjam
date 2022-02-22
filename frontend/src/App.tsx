import React, {useState} from 'react';
import './styles/App.css';
import {SocketProvider} from "./providers/SocketProvider";
import {Router} from "./components/Router";
import {LocalStreamProvider} from "./providers/LocalStreamProvider";


function App() {

    return (
        <SocketProvider>
            <LocalStreamProvider>
                <Router/>
            </LocalStreamProvider>
        </SocketProvider>
    );
}

export default App;
