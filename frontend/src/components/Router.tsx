import React from 'react';
import {
    BrowserRouter,
    Routes,
    Route
} from "react-router-dom";
import {Home} from "./Home";
import Help from "./Help";

export function Router() {
    console.log('Router started')
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Help/>}/>
                <Route path="/:myRole" element={<Home/>}/>
            </Routes>
        </BrowserRouter>
    )
}