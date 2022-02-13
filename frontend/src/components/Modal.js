import React from 'react';
import styled from "styled-components";

export default function Modal({children, show}) {
    if (!show) return (<div/>);
    return (
        <Container>
            <Logo/>
            <Wrapper>
                <Content>
                    <br/>
                    {children}
                </Content>
            </Wrapper>
        </Container>
    );
}


const Container = styled.div`
  /* This way it could be display flex or grid or whatever also. */
  display: block;

  /* Probably need media queries here */
  width: 400px;
  max-width: 90%;

  //height: 240px;
  max-height: 100%;

  position: fixed;

  z-index: 100;

  left: 50%;
  top: 50%;

  /* Use this for centering if unknown width/height */
  transform: translate(-50%, -50%);

  /* If known, negative margins are probably better (less chance of blurry text). */
  /* margin: -200px 0 0 -200px; */

  background: white;
  box-shadow: 0 0 60px 10px rgba(0, 0, 0, 0.8);

  border-radius: 10px;
  padding: 4px;

  align-items: center;
  text-align: center;
`;

const Logo = styled.div`
  width: 100%;
  height: 130px;
  position: absolute;
  top: -80px;
  left: 0;
  background-image: url("/spark-logo.png");
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;
`;

const Wrapper = styled.div`
  display: flex;
  width: 100%;
  height: 100%;
  border: 2px solid gray;
  border-radius: 10px;
  justify-content: center;
  flex-direction: column;
  align-items: center;
`;

const Content = styled.div`
  display:inline-flex;
  flex-wrap: wrap;
  flex-direction: column;
  justify-content: center;
  width: 90%;
`;
