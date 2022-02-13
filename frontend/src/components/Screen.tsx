import React from "react";
import styled from "styled-components";

export default function Screen(props: any) {
    let videoCount = props.children.length;
    let containers = [];
    const flexGap = 1;
    let flexRatio = 100 / Math.ceil(Math.sqrt(videoCount));
    let flex = "0 0 calc(" + flexRatio + "% - " + flexGap + "px)";

    for (let i = 0; i < videoCount; i++) {
        containers.push(
            <VideoContainer key={"vc-" + i} className={"video-container"} style={{flex: flex}}>
                {i}
            </VideoContainer>
        )
    }
    return (
        <VideoGroupContainer>
            {containers}
        </VideoGroupContainer>
    )
}

const VideoGroupContainer = styled.div`
  position: relative;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  background-color: black;
  gap: 1px;
  height: 95%;
`;

const VideoContainer = styled.div`
  background: teal;
  border: 2px solid black;
  border-radius: 1em;
  color: white;
`;
