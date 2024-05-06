import { ConnectWallet, useSDK } from "@thirdweb-dev/react";
import { useUser, useAddress, useLogin } from "@thirdweb-dev/react";
import styles from "../../styles/Home.module.css";
import "../../styles/globals.css"
import { useState } from "react";
import { v4 as uuidv4 } from 'uuid';


export function WalletConnect() {
    const address = useAddress()
    const { login } = useLogin()

    const [signature, setSignature] = useState('N/A')
    const [addres, setAddres] = useState('N/A')

    // const message = 'Please sign me';

    const currentDate = new Date();
    const expirationDate = new Date(currentDate.getTime() + 30 * 60000); // 30 minutes in milliseconds
    const notBefore = new Date(currentDate.getTime() + 10 * 60000); // 10 minutes in milliseconds
    const randomUUID = uuidv4();
    console.log("UUID: ", randomUUID)

    const loginOptions = {
        version: "1", // The current version of the message, which MUST be 1 for this specification.
        chainId: "1", // Chain ID to which the session is bound
        nonce: randomUUID, // randomized token typically used to prevent replay attacks
        "Issued At": currentDate.toISOString(),
        "Expiration Time": expirationDate.toISOString(), // When this message expires
        "Not Before": notBefore.toISOString(), // When this message becomes valid
    };

    const sdk = useSDK();


    /*{
      
      localhost:3000 wants you to sign in with your Ethereum account:
      0x0319B28efEeF6131f1bcC833eCAA86E5c9c75867
  
      Please ensure that the domain above matches the URL of the current website.
  
      Version: 1
      Chain ID: 1
      Nonce: f280c5d0-cfbf-4f26-82ad-0bd8d11fe37c
      Issued At: 2024-05-02T10:38:32.543Z
      Expiration Time: 2024-05-02T10:59:24.973Z
      Not Before: 2024-05-02T10:39:24.973Z
  
    }*/

    const objectToString = (obj: any) => {
        return Object.entries(obj)
            .map(([key, value]) => `${key}: ${value}`)
            .join('\n');
    }

    const message = `localhost:3000 wants you to sign in with your Ethereum account:\n${address}\n\nPlease ensure that the domain above matches the URL of the current website.\n\n${objectToString(loginOptions)}`

    const signMessage = async () => {
        const sig = await sdk?.wallet?.sign(message)

        if (!sig) {
            throw new Error('Failed to sign message')
        }

        setSignature(sig)
    }

    const recoverAddress = async () => {
        const add = sdk?.wallet?.recoverAddress(message, signature)

        if (!add) {
            throw new Error('No Address!');
        }

        setAddres(add)
    }


    return (
        <div  >
            <div>
                <h1 className={styles.div}>
                    Welcome to GreatApe wallet connect App
                </h1>

                <div className={styles.connect}>
                    <ConnectWallet
                        auth={{
                            loginOptional: true,
                        }} />
                </div>
                {address && <p>Address is {address}</p>}

                {/* <button style={styles.button} onClick={() => { login() }}>
                    Sign in with Etheruem
                </button> */}

                <div className={styles.div}>
                    <button onClick={signMessage}>Sign message</button>
                    <br />
                    <p className={styles.div}>Signature: {signature}</p>
                </div>
                {/* 
                <button className={styles.div} onClick={recoverAddress}>Recover Wallet Address</button>
                <p className={styles.div}>Recovered Wallet Address: {addres}</p> */}
            </div>

        </div>
    );
};

export default WalletConnect;
