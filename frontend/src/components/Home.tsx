import GetName from "./GetName";
import React, {useEffect, useState} from "react";
import {Socket} from "./Socket";
import {useParams} from "react-router-dom";
import {Broadcast} from "./Broadcast";
import {Audience} from "./Audience";
// import {useParams} from "react-router-dom";

export const Home = () => {
    // const [myName, setMyName] = useState('');
    let {myRole} = useParams()

    // return (
    //     <>
    //         {
    //             myName ?
    //                 <Socket myName={myName}/> :
    //                 <GetName myName={myName} setMyName={setMyName}/>
    //         }
    //     </>
    // )

    // Choose a name for now
    return myRole ==='broadcast' ? <Broadcast myName={'brd'}/> : <Audience myName={'aud'}/>
    // return <Socket myName={myName}/>
}
