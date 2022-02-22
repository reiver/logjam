import React, {useState} from 'react';
import styled from "styled-components";
import Screen from "./Screen";
import BottomSidebar from "./BottomSidebar";
import LeftSidebar from "./LeftSidebar";
import ControlButtons from "./ControlButtons";
import LocalStream from "./LocalStream";

export default function Main({myName, myRole}: any) {

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


    return (
        <Page>
            <Screen>
                <LocalStream/>
                <video id="video 2" autoPlay playsInline muted/>
                <video id="video 3" autoPlay playsInline muted/>
                <video id="video 4" autoPlay playsInline muted/>
                <video id="video 5" autoPlay playsInline muted/>
            {/*    /!*<video id="video 6" autoPlay playsInline muted/>*!/*/}
            {/*    /!*<video id="video 7" autoPlay playsInline muted/>*!/*/}
            {/*    /!*<video id="video 8" autoPlay playsInline muted/>*!/*/}
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
