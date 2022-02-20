import React, {createContext, ReactChild, useEffect, useState} from 'react';
import './styles/App.css';
import Main from "./components/Main";
import GetName from "./components/GetName";
import {SocketProvider} from "./SocketProvider";
import {MyComponent} from "./components/MyComponent";

function App() {
    // const [myName, setMyName] = useState();


    // useEffect(()=>{
    //     console.log('myName is set to:', myName);
    // }, myName);
    return <SocketProvider>{
        <div>
            <MyComponent/>

            {/*{*/}
            {/*    myName ?*/}
            {/*        <Main myName={myName}/> :*/}
            {/*        <GetName myName={myName} setMyName={setMyName}/>*/}
            {/*}*/}
        </div>

    }</SocketProvider>
}

export default App;
