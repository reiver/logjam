import { signal } from '@preact/signals'
const dialogs = signal([])
import { v4 as uuidv4 } from 'uuid'

import { clsx } from 'clsx'
import { Button, Icon, IconButton, Logo, Tooltip } from 'components'
import { useEffect, useRef, useState } from 'preact/compat'
import Close from 'assets/icons/Close.svg?react'
import metaTagsThumbnail from 'assets/images/metatagsLogo.png'
import reset from 'assets/images/Reset.png'
import { RoundButton } from 'components/common/RoundButton'
import imageFolder from 'assets/images/ImageFolder.png'

export const HostDialogPool = () => {
    console.log("Inside HostDialogPool")
    return (
        <>
            {Object.values(dialogs.value).map((dialog) => {
                if (dialog.type === 'meta-image') return <MetaImageDialog {...dialog} />

            })}
        </>
    )
}

export const makeMetaImageDialog = (oldImage, type, message, onOk, onClose, options = {}) => {
    console.log("inside makeMetaImageDialog")
    const id = uuidv4()
    const destroy = () => {
        const dialogsTmp = { ...dialogs.value }
        delete dialogsTmp[id]
        dialogs.value = dialogsTmp
    }

    dialogs.value = {
        ...dialogs.value,
        [id]: {
            id,
            oldImage,
            type,
            message,
            pointer: !!onClose,
            onOk: () => {
                onOk && onOk()
                destroy()
            },
            onClose: async (image) => {

                onClose && onClose(image)
                destroy()
            },
            ...options,
        },
    }

    return id
}

export const MetaImageDialog = ({
    onOk,
    onClose,
    oldImage,
    message: { title },
    okText = 'Choose Picture',
    cancelText = 'Reset',
    showButtons = true,
    className,
}) => {

    const [selectedImage, setSelectedImage] = useState(oldImage);

    const handleImageChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                setSelectedImage(e.target.result);
                console.log("sleectedImage: ", selectedImage)
            };
            reader.readAsDataURL(file);
        }
    };

    const resetThumbnail = () => {
        setSelectedImage(null)
    }

    console.log("Inside MetaImageDialog")

    return (
        <div class="absolute top-0 left-0 w-full h-full">
            <div class="z-10 absolute w-full h-full bg-black bg-opacity-60" />
            <div
                class={clsx(
                    className,
                    'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full sm:max-w-[40%] sm:border dark:border-gray-1 border-gray-0'
                    , 'w-full sm:w-[416px]'
                )}
            >
                <div class="flex justify-center items-center p-5 relative">
                    <span class="dark:text-white text-black text-bold-12">{title}</span>
                    <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={() => { onClose(selectedImage) }} />
                </div>
                <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />

                <img class="w-[384px] h-full max-h-[384px] max-w-[384px] my-6 mx-4 rounded-md" src={selectedImage ? selectedImage : metaTagsThumbnail}>
                </img>

                {showButtons && (
                    <div class="flex items-center justify-center gap-x-24 p-5 pt-0">

                        <RoundButton image={reset} variant="red" onClick={resetThumbnail}>
                            {cancelText}
                        </RoundButton>

                        <input id="imageInput" type="file" accept="image/*" style="display:none;" onChange={handleImageChange} />
                        <RoundButton image={imageFolder} variant="solid" onClick={() => document.getElementById('imageInput').click()}>
                            {okText}
                        </RoundButton>

                    </div>
                )}
            </div>
        </div>
    )
}

export const HostToastProvider = () => {
    return (
        <div id="toast-provider">
            <HostDialogPool />
        </div>
    )
}