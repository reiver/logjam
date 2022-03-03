import React from 'react';
import './styles/App.css';
import {Router} from "./components/Router";
import {SocketProvider} from "./providers/SocketProvider";
import {PeerConnectionMapProvider} from "./providers/PeerConnectionMapProvider";
import {StreamMapProvider} from "./providers/StreamMapProvider";


export default function App() {
    console.log('APP started')
    return (
        <SocketProvider>
            <StreamMapProvider>
                <PeerConnectionMapProvider>
                    <Router/>
                </PeerConnectionMapProvider>
            </StreamMapProvider>
        </SocketProvider>
    );
}
