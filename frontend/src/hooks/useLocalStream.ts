import {useContext} from "react";
import {LocalStreamContext} from "../providers/LocalStreamProvider";

export const useLocalStream = () => {
    console.log('useLocalStream');
    return useContext(LocalStreamContext);
}
