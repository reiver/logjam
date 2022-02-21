import React from 'react';
import './styles/App.css';
import {SocketProvider} from "./SocketProvider";
import {Home} from "./components/Home";


function App() {
    return (
        <SocketProvider>
            <Home/>
        </SocketProvider>
    );
}

export default App;
