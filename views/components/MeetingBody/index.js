import {
    Attendees,
    AttendeesBottomSheet,
    BottomBarBottomSheet,
    Container,
} from 'components';
import { html } from 'htm';
import { Stage, streamers } from './Stage.js';

export { streamers };

export const MeetingBody = () => {
    return html`<${Container}
            class="relative h-full flex-grow overflow-hidden py-4 flex items-center"
            id="meeting-body"
        >
            <${Stage} />
            <$${Attendees} />
        <//>
        <div class="sm:hidden">
            <${BottomBarBottomSheet} />
            <${AttendeesBottomSheet} />
        </div>`;
};
