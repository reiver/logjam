import { Button, Container, Logo } from 'components';
import { html } from 'htm';
import { leaveMeeting } from '../../pages/meeting.js';

export const TopBar = () => {
    const handleLeaveMeeting = leaveMeeting;
    return html`<div class="w-full bg-white dark:bg-black py-3" id="top-bar">
        <${Container}>
            <div class="grid grid-cols-12">
                <div class="col-span-3"><${Logo} /></div>
                <div class="col-span-6 flex items-center justify-center">
                    <span
                        class="text-black dark:text-white text-center text-bold-14"
                        >Is Your Future Distributed? Welcome to the
                        Fediverse!</span
                    >
                </div>
                <div class="col-span-3 text-right">
                    <${Button} variant="red" onClick=${handleLeaveMeeting}
                        >Leave<//
                    >
                </div>
            <//>
        <//>
    <//>`;
};

export default TopBar;
