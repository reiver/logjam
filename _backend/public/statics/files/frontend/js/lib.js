const GroupVideo = {
    async requestPermissions() {
        try {
            await navigator.mediaDevices.getUserMedia({
                audio: true,
                video: true,
            });
            return true;
        } catch (err) {
            return false;
        }
    }
};