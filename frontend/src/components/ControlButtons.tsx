import React from 'react';
import styled from "styled-components";

const CAMERA_ON = "/cam-on.png";
const CAMERA_OFF = "/cam-off.png";
const MIC_ON = "/mic-on.png";
const MIC_OFF = "/mic-off.png";
const SPARK_LOGO = "/spark-logo.png";

export default function ControlButtons({camera, mic, onMicButtonClick, onCameraClick}: {
    camera: boolean,
    mic: boolean,
    onMicButtonClick: ()=>void,
    onCameraClick: ()=> void
}) {
    return (
        <Container>
            <Image src={camera ? CAMERA_ON : CAMERA_OFF} alt={"icon"} onClick={onCameraClick}/>
            <Image src={SPARK_LOGO} alt={"logo"}/>
            <Image src={mic ? MIC_ON : MIC_OFF} alt={"icon"} onClick={onMicButtonClick}/>
        </Container>
    )
}

const Container = styled.div`
  display: flex;
  flex-direction: row;
  gap: 30px;
  padding-top: 50px;
  height: 100%;
  position: relative;
  justify-content: center;
  text-align: center;
`;

const Image = styled.img`
  height: 25%;
  cursor: pointer;
  &:hover{
    filter: brightness(0) invert(1);
  }
`;