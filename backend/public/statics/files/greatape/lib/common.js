import { SparkRTC } from './spark-rtc.js';

// TODO: set base url
export function getWsUrl(host = null) {
    let baseUrl = null;

    if (host) {
        baseUrl = host;
    } else {
        baseUrl = window.location.href.split('//')[1].split('/')[0];
    }

    const protocol =
        window.location.href.split('//')[0] === 'http:' ? 'ws' : 'wss';
    return `${protocol}://${baseUrl}/ws`;
}

export function createSparkRTC(role, options) {
    // TODO: set role
    if (role === 'broadcast') {
        return createBroadcastSpartRTC(role, options);
    } else {
        return createAudienceSpartRTC(role, options);
    }
}

export const createBroadcastSpartRTC = (role, props) => {
    return new SparkRTC(role, {
        localStreamChangeCallback: props.localStreamChangeCallback,
        remoteStreamCallback: props.remoteStreamCallback,
        remoteStreamDCCallback: props.remoteStreamDCCallback,
        raiseHandConfirmation: props.onRaiseHand,
        userListUpdated: props.onUserListUpdate,
        startProcedure: props.onStart,
        signalingDisconnectedCallback: props.signalingDisconnectedCallback,
        treeCallback: props.treeCallback,
        constraintResults: props.constraintResults,
        updateStatus: props.updateStatus,
        userInitialized: props.onUserInitialized,
    });
};

export const createAudienceSpartRTC = (role, props) => {
    return new SparkRTC(role, {
        remoteStreamCallback: props.remoteStreamCallback,
        remoteStreamDCCallback: props.remoteStreamDCCallback,
        startProcedure: props.onStart,
        userListUpdated: props.onUserListUpdate,
        altBroadcastApprove: props.altBroadcastApprove,
        maxLimitReached: props.maxLimitReached,
        disableBroadcasting: props.disableBroadcasting,
        updateStatus: props.updateStatus,
        userInitialized: props.onUserInitialized,
    });
};
