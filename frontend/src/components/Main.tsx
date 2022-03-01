import React, {useState} from 'react';
import styled from "styled-components";
import Screen from "./Screen";
import BottomSidebar from "./BottomSidebar";
import LeftSidebar from "./LeftSidebar";
import ControlButtons from "./ControlButtons";
import Stream from "./Stream";
import {useStreamMap} from "../hooks/useStreamMap";

export default function Main({myName}: any) {
    console.log('[Render] Main');
    let {streamMap, setStreamMap, enableLocalStream} = useStreamMap();

    const [mic, setMic] = useState(true);
    const [camera, setCamera] = useState(true);

    function handleMicButtonClicked() {
        setMic(!mic);
    }

    function handleCameraButtonClicked() {
        setCamera(!camera);
    }

    const controlButtons = (
        <ControlButtons camera={camera} mic={mic}
                        onMicButtonClick={handleMicButtonClicked}
                        onCameraClick={handleCameraButtonClicked}
        />
    );

    // {myRole === 'broadcast' ? <LocalStream/> : null}

    return (
        <Page>
            <Screen>
                <Stream streamId={'localStream'} mediaStream={streamMap.get('localStream')}/>
                <Stream streamId={'remoteStream'} mediaStream={streamMap.get('remoteStream')}/>
                <video id="video 2" autoPlay playsInline muted/>
                <video id="video 3" autoPlay playsInline muted/>
                <video id="video 4" autoPlay playsInline muted/>
            </Screen>
            <BottomSidebar controlButtons={controlButtons} myName={myName}/>
            <LeftSidebar/>
        </Page>
    )
}

const Page = styled.div`
  background-color: black;
  height: 100vh;
  padding: 1em;
`;
