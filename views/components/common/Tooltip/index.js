import { cloneElement, toChildArray } from 'https://esm.sh/preact';
import { useEffect, useRef } from 'preact';
import { deviceSize } from '../../MeetingBody/Stage.js';

export const Tooltip = ({ children, label }) => {
    const ref = useRef();
    const component = toChildArray(children)[0];
    const tippyInstance = useRef();
    useEffect(() => {
        if ((ref.current.base || ref.current) && deviceSize.value !== 'xs')
            tippyInstance.current = tippy(ref.current.base || ref.current, {
                content: label,
                arrow: false,
                hideOnClick: true,
            });

        return () => {
            if (tippyInstance.current) {
                tippyInstance.current.destroy();
            }
        };
    }, [deviceSize.value]);

    useEffect(() => {
        if (tippyInstance.current) tippyInstance.current.setContent(label);
    }, [label]);

    if (component) return cloneElement(component, { ref });
    return null;
};
