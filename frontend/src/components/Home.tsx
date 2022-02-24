import GetName from "./GetName";
import {useState} from "react";
import {Socket} from "./Socket";

export const Home = () => {
    const [myName, setMyName] = useState();

    return (
        <>
            {
                myName ?
                    <Socket myName={myName}/> :
                    <GetName myName={myName} setMyName={setMyName}/>
            }
        </>
    )
}
