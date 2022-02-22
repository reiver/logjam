import React from 'react';
import './styles/App.css';
import {SocketProvider} from "./providers/SocketProvider";
import {Router} from "./components/Router";


function App() {
    return (
        <SocketProvider>
            <Router/>
        </SocketProvider>
    );
}

export default App;
