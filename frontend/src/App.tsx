import React from 'react';
import './styles/App.css';
import {Router} from "./components/Router";
import {SocketProvider} from "./providers/SocketProvider";
import {LocalStreamProvider} from "./providers/LocalStreamProvider";
import {PeerConnectionArrayProvider} from "./providers/PeerConnectionArrayProvider";


export default function App() {

    return (
        <SocketProvider>
            <LocalStreamProvider>
                <PeerConnectionArrayProvider>
                    <Router/>
                </PeerConnectionArrayProvider>
            </LocalStreamProvider>
        </SocketProvider>
    );
}
