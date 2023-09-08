import { signal } from '@preact/signals';
import 'external-svg-loader';
import { html, render } from 'htm';
import App from './App.js';

export const userInteractedWithDom = signal(false);

const events = [
    'mousedown',
    'click',
    'mouseup',
    'keydown',
    'keypress',
    'scroll',
    'touchstart',
    'touchend',
    'touchmove',
];
const onInteract = () => {
    userInteractedWithDom.value = true;
    events.forEach((event) => {
        document.removeEventListener(event, onInteract);
    });
    if (interval) clearInterval(interval);
};

const interval = setInterval(() => {
    if (!navigator.userActivation) {
        events.forEach((event) => {
            document.addEventListener(event, onInteract);
        });
        clearInterval(interval);
        return;
    }
    if (
        navigator.userActivation.isActive &&
        navigator.userActivation.hasBeenActive
    ) {
        onInteract();
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
