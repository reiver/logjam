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
import { css } from '@emotion/react'
import { CSSData, PocketBaseManager } from 'lib/helperAPI'
const pbApi = new PocketBaseManager()

export const HostDialogPool = () => {
    console.log("Inside HostDialogPool")
    return (
        <>
            {Object.values(dialogs.value).map((dialog) => {
                if (dialog.type === 'meta-image') return <MetaImageDialog {...dialog} />
                else if (dialog.type === 'css-files') return <CssFilesDialog {...dialog} />

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
            onClose: async (image, imageFile) => {

                onClose && onClose(image, imageFile)
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
    const [imageFile, setImageFile] = useState(null)

    const handleImageChange = (event) => {
        const file = event.target.files[0];
        setImageFile(file)

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
                    <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={() => { onClose(selectedImage, imageFile) }} />
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

export const makeCssFilesDialog = (cssFiles, hostId, oldIndex, type, message, onOk, onClose, options = {}) => {
    console.log("inside makeCssFilesDialog")
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
            cssFiles,
            hostId,
            oldIndex,
            type,
            message,
            pointer: !!onClose,
            onOk: () => {
                onOk && onOk()
                destroy()
            },
            onClose: async (cssFile, index) => {

                onClose && onClose(cssFile, index)
                destroy()
            },
            ...options,
        },
    }

    return id
}

const createNewCSS = async (cssData) => {
    var newCSS = await pbApi.createCSS(cssData)
    console.log("new CSS Created: ", newCSS)
    return newCSS;
}

export const CssFilesDialog = ({
    onOk,
    onClose,
    cssFiles,
    hostId,
    oldIndex,
    message: { title },
    okText = 'Choose Picture',
    cancelText = 'Reset',
    showButtons = true,
    className,
}) => {
    let selectedFileIndex = oldIndex
    let customStyles = null
    let uploadedFile = null

    //set default selection
    setTimeout(() => {
        const radioInput = document.getElementById(`file${selectedFileIndex}`) as HTMLInputElement
        radioInput.checked = selectedFileIndex === selectedFileIndex
        radioInput.style.accentColor = 'black'

    }, 50);

    console.log("Inside MetaImageDialog")

    const handleFileClick = (index = -1, vanish = true) => {

        // Now, the 'index' variable will contain the index of the matching device (or -1 if none found).

        // Check if the clicked device is already selected
        if (selectedFileIndex === index) {
            // If it's already selected, deselect it by setting the selectedDeviceIndex to -1
            selectedFileIndex = -1
        } else {
            // If it's not selected, select it by setting the selectedDeviceIndex to the clicked index
            selectedFileIndex = index
        }

        console.log('handleDeviceClick: ', index)
        // Update the radio button's checked attribute based on the selectedDeviceIndex
        const radioInput = document.getElementById(`file${index}`) as HTMLInputElement
        console.log('radioInput: ', radioInput)
        if (radioInput) {
            radioInput.checked = selectedFileIndex === index
            radioInput.style.accentColor = 'black'
            if (vanish) {
                setTimeout(() => {
                    console.log('Selected file: ', cssFiles[selectedFileIndex])
                    onClose(cssFiles[selectedFileIndex], selectedFileIndex)
                }, 200)
            }
        }
    }

    const handleCssFileUpload = async (event) => {

        const fileInput = event.target;
        const fileLabel = document.getElementById('fileLabel');
        uploadedFile = fileInput.files[0]

        // Check if files were selected
        if (fileInput.files.length > 0) {
            // Update label with the first selected file's name
            fileLabel.textContent = fileInput.files[0].name;
        } else {
            // No file selected, reset label text
            fileLabel.textContent = "Choose CSS file";
        }


        setCustomCssContent(event, (content) => {

            // Regular expression to match class names
            const cssClassRegex = /\.([a-zA-Z0-9_-]+)/g;

            // Match all class names in the CSS content and add them to the array
            let match;
            while ((match = cssClassRegex.exec(content)) !== null) {
                if (!cssClassNames.includes(match[1])) {
                    cssClassNames.push(match[1]);
                }
            }
            console.log("Css Class names: ", cssClassNames)

            if (isValidCSS(content)) {
                console.log("Is Valid true")
                if (content) {

                    // cssFiles.add(fileInput.files[0])

                    customStyles = content


                    // //user uploaded valid css... Now save this css to DB
                    var cssData = new CSSData('', fileInput.files[0].name, customStyles, hostId)
                    createNewCSS(cssData)


                    return
                }
            } else {
                console.log("Is Valid false")
            }
        });
    };

    const setCustomCssContent = (event, setContentCallback) => {
        const file = event.target.files[0];


        if (file) {
            const reader = new FileReader();

            reader.onload = (e) => {
                const content = e.target.result;
                setContentCallback(content);
            };

            reader.readAsText(file);
        } else {
            // Handle the case where no file is selected
            setContentCallback(null);
        }
    };

    // Define an array to store all the CSS class names from the provided CSS content
    const cssClassNames: string[] = [];

    // Check if all required classes are present in the CSS content
    const requiredClasses = [
        'greatape-stage-host',
        'greatape-stage-host-audience-1',
        'greatape-stage-host-screenshare',
        'greatape-stage-host-screenshare-audience-1',
        'greatape-stage-host-audience-2',
        'greatape-stage-host-audience-3',
        'greatape-gap-in-videos',
        'greatape-host-video',
        'greatape-share-screen-video',
        'greatape-audience-video',
        'greatape-video-name',
        'greatape-video-name-background',
        'greatape-attendees-list',
        'greatape-attendees-count',
        'greatape-attendees-item',
        'greatape-attendees-item-role',
        'greatape-meeting-link',
        'greatape-meeting-link-background'
    ];

    function isValidCSS(cssContent: string): boolean {

        const allClassesPresent = requiredClasses.every(className => cssClassNames.includes(className));

        // Regular expression to match CSS rules
        const cssRuleRegex = /[^{]*\{[^}]*\}/g;

        // Match all CSS rules in the content
        const matches = cssContent.match(cssRuleRegex);

        // If matches are found and every match has a valid structure, and all required classes are present, return true
        return (
            matches !== null &&
            matches.every(match => isValidCSSRule(match)) &&
            allClassesPresent
        );
    }

    function isValidCSSRule(cssRule: string): boolean {
        // Regular expression to match a single CSS rule
        const cssRuleStructureRegex = /^\s*([^\{\}]+)\s*\{([^\{\}]*)\}\s*$/;

        // Check if the CSS rule matches the expected structure
        return cssRuleStructureRegex.test(cssRule);
    }

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
                    <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={() => {
                        if (selectedFileIndex == -1 && uploadedFile != null) {
                            onClose(uploadedFile, 0)
                        } else {
                            if (selectedFileIndex != -1) {
                                onClose(cssFiles[selectedFileIndex], selectedFileIndex)
                            } else {
                                onClose(null, selectedFileIndex)
                            }
                        }
                    }} />
                </div>
                <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />

                <div class={clsx(
                    "overflow-y-auto", `${cssFiles != null && cssFiles.length > 0 ? "h-64" : ""}`
                )}>
                    <style>
                        {`
                            /* Customizing scrollbar styles for WebKit browsers */
                            ::-webkit-scrollbar {
                            width: 3px; /* Adjust the width as needed */
                            }

                            ::-webkit-scrollbar-thumb {
                            background-color: #A8A8A8; /* Adjust the color as needed */
                            }

                            ::-webkit-scrollbar-track {
                            background-color: #EBEBEB; /* Adjust the color as needed */
                            }
                        `}
                    </style>

                    <div class="sm:py-4 py-2 sm:px-6 px-4 flex items-center justify-between relative">
                        <div class={clsx('text-bold-12 text-gray-3')}>Upload New</div> <label id="fileLabel" for="cssFileInput" class={clsx('text-bold-12 text-gray-1 cursor-pointer')}>
                            Choose CSS file
                        </label>
                        <input id="cssFileInput" type="file" class="hidden" onChange={(event) => { handleCssFileUpload(event) }} />

                    </div>
                    <hr class="dark:border-gray-2 border-gray-0 mx-4 my-1 sm:mx-6 sm:my-1" />

                    <form>
                        <div class="sm:pb-4 pb-2">

                            <div class="sm:py-4 py-2 sm:px-6 px-4 rounded-md flex items-center cursor-pointer" onClick={() => handleFileClick(-1)}>

                                <div class="text-left text-bold-12 flex-1">Default</div>
                                <label class="flex items-right flex-0">
                                    <input type="radio" name="devices" checked={true} id={`file${-1}`} />
                                </label>
                            </div>

                            {cssFiles != null && cssFiles.length > 0 && cssFiles.map((file, index) => (

                                <div>
                                    <hr class="dark:border-gray-2 border-gray-0 mx-4 my-1 sm:mx-6 sm:my-1" />
                                    <div class="sm:py-4 py-2 sm:px-6 px-4 rounded-md flex items-center cursor-pointer" onClick={() => handleFileClick(index)}>
                                        <div class="text-left text-bold-12 flex-1">{file.name}</div>
                                        <label class="flex items-right flex-0">
                                            <input type="radio" name="devices" id={`file${index}`} />
                                        </label>
                                    </div>
                                </div>


                            ))}
                        </div>
                    </form>

                </div>

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