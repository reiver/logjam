import React from 'react';
import styled from "styled-components";

export default function BottomSidebar({controlButtons, myName}: any){

    return(
        <Wrapper>
            <Slider/>
            <Container >
                <Welcome>{myName}</Welcome>
                {controlButtons}
            </Container>
        </Wrapper>
    )
}


const Welcome = styled.p`
  color: white;
  margin-left: 20px;
  position: relative;
  top: -30px;
`;


const Wrapper = styled.div`
  bottom: 0;
  left: 0;
  position: fixed;
  width: 100%;
  height: 30px;
  //overflow: auto;
  margin: 0 auto;
  -webkit-transition: all 1s ease;
  -moz-transition: all 1s ease;
  -o-transition: all 1s ease;
  transition: all 1s ease;
  color: white;

  &:hover {
    -webkit-transition: all 1s ease;
    -moz-transition: all 1s ease;
    -o-transition: all 1s ease;
    transition: all 1s ease;
    height: 50%;
  }
`;

const Container = styled.div`
  width: 100%;
  height: 100%;
  position: relative;
  top: 0;
  left: 0;
  background-color: rgba(36, 46, 66, .5);
`;

const Slider = styled.div`
  width: 20%;
  height: 4px;
  border-radius: 10px;
  position: absolute;
  left: 40%;
  bottom: 0;
  background-color: #26a5f6;
  z-index: 1;
`;
