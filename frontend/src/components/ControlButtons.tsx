import React from 'react';
import styled from "styled-components";

const CAMERA_ON = "/cam-on.png";
const CAMERA_OFF = "/cam-off.png";
const MIC_ON = "/mic-on.png";
const MIC_OFF = "/mic-off.png";

export default function ControlButtons({camera, mic}: {
    camera: boolean,
    mic: boolean
}) {
    return (
        <Wrapper>
            <Img src={camera ? CAMERA_ON :  CAMERA_OFF} alt={"icon"}/>
            <Img src={mic ? MIC_ON :  MIC_OFF} alt={"icon"}/>
        </Wrapper>
    )
}

const Wrapper = styled.div`
  //position: absolute;
  //bottom: 0;
  display: flex;
  justify-content: center;
  //z-index: 2;
  padding: 1%;
  //height: 10vh;
  //width: 100vw;
  //background-color: rgba(50, 82, 123, .5);
`;

const Img = styled.img`
  //position: relative;
  //border: 2px solid white;
  height: 5%;
  width: 5%;
  margin: auto 1%;
`;
