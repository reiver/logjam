import React from 'react';
import styled from "styled-components";

export default function Logo(){
    return(
        <Container/>
    )
}

const Container = styled.div`
  width: 100%;
  height: 120px;
  position: absolute;
  top: -80px;
  left: 0;
  background-image: url("/spark-logo.png");
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;
`;
