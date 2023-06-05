import { cloneElement, toChildArray } from 'https://esm.sh/preact';
import { useEffect, useRef } from 'preact';

export const Tooltip = ({ children, label }) => {
    const ref = useRef();
    const component = toChildArray(children)[0];
    const tippyInstance = useRef();
    useEffect(() => {
        if (ref.current.base || ref.current)
            tippyInstance.current = tippy(ref.current.base || ref.current, {
                content: label,
                arrow: false,
                hideOnClick: false,
            });
    }, []);

    useEffect(() => {
        tippyInstance.current.setContent(label);
    }, [label]);

    if (component) return cloneElement(component, { ref });
    return null;
};
