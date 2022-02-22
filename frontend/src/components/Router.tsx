import React from 'react';
import {
    BrowserRouter,
    Routes,
    Route
} from "react-router-dom";
import {Home} from "./Home";

export function Router() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/:myRole" element={<Home/>}/>
            </Routes>
        </BrowserRouter>
    )
}