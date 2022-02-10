import React from 'react';
import styled from "styled-components";
import Screen from "./Screen";
import ControlButtons from "./ControlButtons";

export default function Main() {
    return (
        <Page>
            <Screen>
                <video id="video 1" autoPlay playsInline muted/>
                <video id="video 2" autoPlay playsInline muted/>
                <video id="video 3" autoPlay playsInline muted/>
                <video id="video 4" autoPlay playsInline muted/>
                <video id="video 5" autoPlay playsInline muted/>
                <video id="video 6" autoPlay playsInline muted/>
                <video id="video 7" autoPlay playsInline muted/>
                <video id="video 8" autoPlay playsInline muted/>
            </Screen>
            <Footer>
                <FooterButton/>
                <div id="container">
                    <ControlButtons camera={true} mic={false}/>
                </div>
            </Footer>
        </Page>
    )
}

const Page = styled.div`
  background-color: black;
  border: 5px solid #16253f;
  height: 100vh;
  padding: 1em;
`;

const Footer = styled.div`
  //background-color: red;
  bottom: 0;
  left: 0;
  position: fixed;
  width: 100%;
  height: 1%;
  overflow: hidden;
  margin: 0 auto;
  -webkit-transition: all 1s ease;
  -moz-transition: all 1s ease;
  -o-transition: all 1s ease;
  transition: all 1s ease;
  z-index: 999;

  &:hover {
    -webkit-transition: all 1s ease;
    -moz-transition: all 1s ease;
    -o-transition: all 1s ease;
    transition: all 1s ease;
    height: 20%;
  }

  #container {
    margin-top: 5px;
    width: 100%;
    height: 100%;
    position: relative;
    top: 0;
    left: 0;
    background: rgba(36, 46, 66, .3);
  }
`;

const FooterButton = styled.div`
  width: 50px;
  height: 1px;
  border: #727172 1px solid;
  border-radius: 1px;
  margin: 0 auto;
  position: relative;
  -webkit-transition: all 1s ease;
  -moz-transition: all 1s ease;
  -o-transition: all 1s ease;
  transition: all 1s ease;

  &:hover {
    width: 35px;
    height: 35px;
    border: #3A3A3A 12px solid;
    -webkit-transition: all 1s ease;
    -moz-transition: all 1s ease;
    -o-transition: all 1s ease;
    transition: all 1s ease;
    position: relative;
  }
`;
