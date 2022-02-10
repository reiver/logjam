import React from 'react';
import './styles/App.css';
import Main from "./components/Main";
import GraphVisualization from "./components/GraphVisualization";
import data from "./data";

function App() {
    return (
        <div className="App">
            {/*<Main/>*/}
            <GraphVisualization data={data}/>
        </div>
    );
}

export default App;
