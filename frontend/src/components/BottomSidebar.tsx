import React from 'react';
import styled from "styled-components";

export default function BottomSidebar({controlButtons}: any){

    return(
        <Wrapper>
            <Slider/>
            <Container >
                {controlButtons}
            </Container>
        </Wrapper>
    )
}

const Wrapper = styled.div`
  bottom: 0;
  left: 0;
  position: fixed;
  width: 100%;
  height: 3px;
  //overflow: auto;
  margin: 0 auto;
  -webkit-transition: all 1s ease;
  -moz-transition: all 1s ease;
  -o-transition: all 1s ease;
  transition: all 1s ease;

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
  height: 3px;
  border: 1px solid #5691f8;
  border-radius: 10px;
  position: absolute;
  left: 40%;
  bottom: 0;
  background-color: #5691f8;
`;
