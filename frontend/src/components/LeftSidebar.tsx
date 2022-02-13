import React from 'react';
import styled from "styled-components";
import GraphVisualization from "./GraphVisualization";
import data from "../data";

export default function LeftSidebar() {
    return (
        <Wrapper>
            <Slider/>
            <Container>
                <GraphVisualization data={data}/>
            </Container>
        </Wrapper>
    )
}

const Wrapper = styled.div`
  top: 0;
  left: 0;
  position: fixed;
  width: 3px;
  height: 100%;
  overflow: auto;
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
    width: 50%;
  }
`;

const Container = styled.div`
  width: 100%;
  height: 100%;
  position: relative;
  top: 0;
  left: 0;
  background: rgba(36, 46, 66, .3);
`;

const Slider = styled.div`
  width: 4px;
  height: 40%;
  border-radius: 10px;
  position: absolute;
  left: 0;
  top: 30%;
  background-color: #26a5f6;
`;
