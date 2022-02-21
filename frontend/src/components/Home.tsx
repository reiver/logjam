import GetName from "./GetName";
import {useState} from "react";
import {Socket} from "./Socket";
import {useParams} from "react-router-dom";

export const Home = () => {
    let { myRole } = useParams();
    const [myName, setMyName] = useState();

    return (
        <>
            {
                myName ?
                    <Socket myName={myName} myRole={myRole}/> :
                    <GetName myName={myName} setMyName={setMyName}/>
            }
        </>
    )
}
