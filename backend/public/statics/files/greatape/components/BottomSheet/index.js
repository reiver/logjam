import { Icon } from 'components';
import { html } from 'htm';

export const BottomSheet = ({ open, children, onClose, title }) => {
    if (!open) return null;
    return html`<div
        class="min-w-[100vw] min-h-[100vh] absolute top-0 left-0 dark:text-white-f-9"
    >
        <div
            class="overlay bg-black bg-opacity-50 w-full min-h-[100vh]"
            onClick=${onClose}
        />
        <div
            class="bottom-0 to-top absolute max-h-[80%] w-full dark:bg-gray-3 bg-white-f flex flex-col rounded-t-[16px]"
        >
            <div class="flex justify-between relative px-4 py-4">
                <${Icon}
                    icon="Close"
                    onClick=${onClose}
                    class="cursor-pointer"
                />
                <span
                    class="dark:text-white absolute left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2 text-bold-12"
                    >${title}</span
                >
            </div>
            <div class="px-4 overflow-auto pb-1">
                <div>${children}</div>
            </div>
        </div>
    </div>`;
};