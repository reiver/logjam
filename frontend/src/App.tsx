import React from 'react';
import './styles/App.css';
import {Router} from "./components/Router";
import {SocketProvider} from "./providers/SocketProvider";
import {LocalStreamProvider} from "./providers/LocalStreamProvider";
import {PeerConnectionMapProvider} from "./providers/PeerConnectionMapProvider";


export default function App() {

    return (
        <SocketProvider>
            <LocalStreamProvider>
                <PeerConnectionMapProvider>
                    <Router/>
                </PeerConnectionMapProvider>
            </LocalStreamProvider>
        </SocketProvider>
    );
}
