import GetName from "./GetName";
import React, {useEffect, useState} from "react";
import {Socket} from "./Socket";
// import {useParams} from "react-router-dom";

export const Home = () => {
    const [myName, setMyName] = useState('');
    // let {myRole} = useParams()

    return (
        <>
            {
                myName ?
                    <Socket myName={myName}/> :
                    <GetName myName={myName} setMyName={setMyName}/>
            }
        </>
    )

    // Choose a name for now
    // return <Socket myName={myRole + '-user'}/>
    // return <Socket myName={myName}/>
}
