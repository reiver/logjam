import React, {useEffect, useState} from "react";
import {makeId} from "../helpers/uid";
import Modal from "./Modal";
import '../styles/Form.css';

export default function GetName({setMyName}: any) {
    const [input, setInput] = useState('');
    const [showModal, setShowModal] = useState(true);

    useEffect(() => {
        let savedName = localStorage.getItem('logjam-myName');
        if (savedName) {
            setInput(savedName);
        }
    }, []);

    const handleSubmit = async (e: any) => {
        e.preventDefault();
        let name = input || makeId(20);
        // store the user in localStorage
        try {
            localStorage.setItem('logjam-myName', name);
        } catch (e) {
            console.log(e)
        }
        setShowModal(false);
        setMyName(name);
    }

    return (
        <Modal show={showModal}>
            <h2>Welcome to Group Video!</h2>
            <form className={"form"} onSubmit={handleSubmit}>
                <label htmlFor="username" style={{height: "2.5em", color: "black"}}>Name: </label>
                <input
                    type="text"
                    value={input}
                    style={{height: "2.5em"}}
                    placeholder="enter a nickname"
                    onChange={({target}) => setInput(target.value)}
                />
                <button style={{height: "2.5em"}} type="submit">Enter</button>
            </form>
        </Modal>
    )
}

