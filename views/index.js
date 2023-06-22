import { signal } from '@preact/signals';
import 'external-svg-loader';
import { html, render } from 'htm';
import App from './App.js';

export const userInteractedWithDom = signal(false);

const interval = setInterval(() => {
    if (
        navigator.userActivation.isActive &&
        navigator.userActivation.hasBeenActive
    ) {
        userInteractedWithDom.value = true;
        clearInterval(interval);
    }
}, 500);

render(html`<${App} />`, document.getElementById('root'));

const documentHeight = () => {
    const doc = document.documentElement;
    doc.style.setProperty('--doc-height', `${window.innerHeight}px`);
};
window.addEventListener('resize', documentHeight);
documentHeight();
setTimeout(documentHeight, 300);
