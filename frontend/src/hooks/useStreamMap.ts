import {useContext} from "react";
import {StreamMapContext} from "../providers/StreamMapProvider";

export const useStreamMap = () => {
    return useContext(StreamMapContext);
}
