import React, { useState, useEffect } from 'react';
import { ConnectWallet, useSDK, useAddress, useLogin } from "@thirdweb-dev/react";
import styles from "../../styles/Home.module.css";
import "../../styles/globals.css";
import { v4 as uuidv4 } from 'uuid';
import { lookUpENS, resolveENS } from "./resolveens";
import { PocketBaseManager, walletData } from "lib/helperAPI";

const pbApi = new PocketBaseManager();

const saveWalletData = async (data) => {
    var wallet = await pbApi.saveWalletAddress(data)
    console.log("Wallet Saved: ", wallet)
}

const WalletConnect = () => {
    const address = useAddress();
    const { login } = useLogin();
    const sdk = useSDK();

    const [signature, setSignature] = useState('N/A');
    const [message, setMessage] = useState('N/A');
    const [resolvedEnsAddress, setResolvedEnsAddress] = useState('N/A');
    const [resolvedWalletAddress, setResolvedWalletAddress] = useState('N/A');

    const ens = "vitalik.eth";
    const [ensInput, setEnsInput] = useState(''); // State to store the user's ENS input

    const handleEnsInputChange = (event) => {
        setEnsInput(event.target.value); // Update the ENS input state when the user types
    };

    useEffect(() => {
        const currentDate = new Date();
        const expirationDate = new Date(currentDate.getTime() + 30 * 60000);
        const notBefore = new Date(currentDate.getTime() + 10 * 60000);
        const randomUUID = uuidv4();

        const loginOptions = {
            version: "1",
            chainId: "1",
            nonce: randomUUID,
            "Issued At": currentDate.toISOString(),
            "Expiration Time": expirationDate.toISOString(),
            "Not Before": notBefore.toISOString(),
        };

        const objectToString = (obj) => {
            return Object.entries(obj)
                .map(([key, value]) => `${key}: ${value}`)
                .join('\n');
        }

        const messageText = `${window.location.host} wants you to sign in with your Ethereum account:\n${address}\n\nPlease ensure that the domain above matches the URL of the current website.\n\n${objectToString(loginOptions)}`;

        setMessage(messageText);
    }, [address]);

    const signMessage = async () => {
        const sig = await sdk?.wallet?.sign(message);


        if (!sig) {
            throw new Error('Failed to sign message');
        }

        var data = new walletData(address, message, sig)
        saveWalletData(data)

        setSignature(sig);
    }

    const resolveENSAddress = async () => {
        try {
            const resolvedName = await resolveENS(ensInput);
            if (resolvedName) {
                setResolvedEnsAddress(resolvedName);
            } else {
                setResolvedEnsAddress('ENS name could not be resolved');
            }
        } catch (error) {
            console.error('Error resolving ENS name:', error);
        }
    }

    const lookupAddress = async () => {
        try {
            const resolvedName = await lookUpENS(address);
            if (resolvedName) {
                setResolvedWalletAddress(resolvedName);
            } else {
                setResolvedWalletAddress('Wallet Address could not be resolved')
            }
        } catch (error) {
            console.error('Error resolving Wallet address:', error);
        }
    }

    return (
        <div className={styles.page}>
            <div>
                <h1 className={styles.div}>
                    <b>Welcome to GreatApe wallet connect App</b>
                </h1>

                <div className={styles.connect}>
                    <ConnectWallet
                        auth={{
                            loginOptional: true,
                        }} />
                </div>
                {address && <p>Address is <b>{address}</b></p>}

                <div className={styles.div}>
                    <button onClick={lookupAddress}>Lookup Address</button>
                    <br />
                    <p className={styles.div}>Linked ENS is: <b>{resolvedWalletAddress}</b></p>
                </div>

                <div className={styles.div}>
                    <button onClick={signMessage}>Sign message</button>
                    <br />
                    <p className={styles.div}>Message: <b>{message}</b></p>
                    <p className={styles.div}>Signature: <b>{signature}</b></p>
                </div>


                <div className={styles.div}>
                    <input className={styles.input}
                        type="text"
                        value={ensInput}
                        onChange={handleEnsInputChange}
                        placeholder="Enter ENS"
                    />
                </div>
                <div className={styles.div}>

                    <button onClick={resolveENSAddress}>Resolve ENS</button>
                    <br />
                    <p className={styles.div}>Ens: <b>{ensInput}</b></p>
                    <p className={styles.div}>Resolved Address: <b>{resolvedEnsAddress}</b></p>
                </div>

                {/* 
                <button className={styles.div} onClick={recoverAddress}>Recover Wallet Address</button>
                <p className={styles.div}>Recovered Wallet Address: {addres}</p> */}
            </div>
        </div>
    );
};

export default WalletConnect;
