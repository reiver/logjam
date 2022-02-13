import React from 'react';
import styled from "styled-components";
import Logo from "./Logo";

const CAMERA_ON = "/cam-on.png";
const CAMERA_OFF = "/cam-off.png";
const MIC_ON = "/mic-on.png";
const MIC_OFF = "/mic-off.png";

export default function ControlButtons({camera, mic}: {
    camera: boolean,
    mic: boolean
}) {
    return (
        <Container>
            <Image src={camera ? CAMERA_ON : CAMERA_OFF} alt={"icon"}/>
            <Image src={"/spark-logo.png"} alt={"logo"}/>
            <Image src={mic ? MIC_ON : MIC_OFF} alt={"icon"}/>
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
`;
