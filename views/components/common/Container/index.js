import clsx from 'clsx';
import { html } from 'htm';

export const Container = ({ children, class: className, ...props }) => {
    return html`<div
        class=${clsx('mx-auto w-full px-4 ', className)}
        ...${props}
    >
        ${children}
    </div>`;
};

export default Container;
