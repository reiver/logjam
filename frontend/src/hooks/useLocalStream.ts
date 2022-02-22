import {useContext} from "react";
import {LocalStreamContext} from "../providers/LocalStreamProvider";

export const useLocalStream = () => {
    return useContext(LocalStreamContext);
}
