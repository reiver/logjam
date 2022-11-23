/**
 * This is a class that is analogous to the DOM's `EventTarget` API.
 *
 * It is the class for adding event listeners, and emitting events.
 *
 * Usage:
 *
 *     const emitter = new EventEmitter()
 *
 *     emitter.addEventListener('foo', (value) => {
 *       // Do whatever with the value.
 *     });
 *
 *     // Then somewhere in your code…
 *     emitter.dispatchEvent('foo', 'the value');
 */
 class EventEmitter {
    /**
     * Initializes a new instance of EventEmitter
     */
    constructor() {
      this._listeners = new Map();
    }
  
    /**
     * Adds an event listener to the labeled event.
     * @param event A string to represent an event name.
     * @param callback A callback that will be called when an event is emitted.
     */
    addEventListener(event, callback) {
      if (!this._listeners.has(event)) {
        this._listeners.set(event, new Set());
      }
      this._listeners.get(event).add(callback);
    }
  
    /**
     * Removes the specified callback, from the list of callbacks.
     * @param event A string to represent an event name.
     * @param callback The callback to remove.
     */
    removeEventListener(event, callback) {
      if (this._listeners.has(event)) {
        const listeners = this._listeners.get(event);
        listeners.delete(callback);
        if (listeners.size <= 0) {
          this._listeners.delete(event);
        }
      }
    }
  
    /**
     * Dispatches an event, and invokes all the callback listening in on the
     * event, as labeled by the `event` parameter.
     * @param event A string that represents which event to emit.
     * @param value The value to emit.
     */
    dispatchEvent(event, value) {
      if (this._listeners.has(event)) {
        for (const listener of this._listeners.get(event)) {
          listener(value);
        }
      }
    }
  }
  
  function sendRequest(iframe, requestType, data, responseType) {
    return new Promise((resolve) => {
      let responded = false;
      let commandId = nextCount();
  
      const sendMessage = () => {
        setTimeout(() => {
          if (responded) {
            return;
          }
          try {
            iframe.contentWindow.postMessage(
              JSON.stringify({
                type: requestType,
                data: { ...data, commandId },
              }),
              new URL(iframe.src).origin
            );
          } catch (e) {}
  
          sendMessage();
        }, 16);
      };
  
      const onMessage = (e) => {
        if (responded) {
          return;
        }
        try {
          const response = JSON.parse(e.data);
          if (response.type === responseType) {
            if (response.data.commandId === commandId) {
              responded = true;
              resolve(response.data);
              window.removeEventListener("message", onMessage);
            }
          }
        } catch (e) {}
      };
  
      window.addEventListener("message", onMessage);
      sendMessage();
    });
  }
  
  let id = 0;
  function nextCount() {
    return id++;
  }
  
  /**
   * The main Group Video class, for Multibrain's seminar video conferencing
   * applications.
   */
  window.GroupVideo = class GroupVideo extends EventEmitter {
    static get PARTICIPANT_ARRIVED() {
      return "PARTICIPANT_ARRIVED";
    }
    static get PARTICIPANT_LEFT() {
      return "PARTICIPANT_LEFT";
    }
    static get PARTICIPANT_NAME_CHANGED() {
      return "PARTICIPANT_NAME_CHANGED";
    }
    static get PARTICIPANT_PROMOTED() {
      return "PARTICIPANT_PROMOTED";
    }
    static get PARTICIPANT_DEMOTED() {
      return "PARTICIPANT_DEMOTED";
    }
    static get MESSAGE_RECEIVED() {
      return "MESSAGE_RECEIVED";
    }
    static get WE_ARE_DEMOTED() {
      return "WE_ARE_DEMOTED";
    }
    static get WE_ARE_PROMOTED() {
      return "WE_ARE_PROMOTED";
    }
  
    /**
     * Initializes a new instance of GroupVideo.
     *
     * The constructor expects an IFrame object that is already navigated to
     * Group Video.
     *
     * Don't know how to create an IFrame for GroupVideo? Then use the
     * `createIFrame` method.
     *
     * The whole purpose of the construtor is to internally add an event listener
     * to the
     * @param {HTMLIFrameElement} iframe The IFrame element that will be hosting
     *   the Group Video instance.
     * @see GroupVideo.createIFrame
     */
    constructor(iframe) {
      super();
  
      this._participants = new Map();
      this._onMessage = ({ data }) => {
        try {
          const { type, data: body } = JSON.parse(data);
          switch (type) {
            case GroupVideo.PARTICIPANT_ARRIVED:
              this._participants.set(body.id, body);
              this.dispatchEvent(GroupVideo.PARTICIPANT_ARRIVED, body);
              break;
            case GroupVideo.PARTICIPANT_LEFT:
              this._participants.delete(body.id);
              this.dispatchEvent(GroupVideo.PARTICIPANT_LEFT, body);
              break;
            case GroupVideo.PARTICIPANT_NAME_CHANGED:
              {
                const participant = this._participants.get(body.id);
                if (participant) {
                  participant.name = body.name;
                  this.dispatchEvent(GroupVideo.PARTICIPANT_NAME_CHANGED, body);
                }
              }
              break;
            case GroupVideo.PARTICIPANT_PROMOTED:
              {
                const participant = this._participants.get(body);
                if (participant) {
                  participant.isParticipating = true;
                  this.dispatchEvent(GroupVideo.PARTICIPANT_PROMOTED, {
                    id: body,
                  });
                }
              }
              break;
            case GroupVideo.PARTICIPANT_DEMOTED:
              {
                const participant = this._participants.get(body);
                if (participant) {
                  participant.isParticipating = false;
                  this.dispatchEvent(GroupVideo.PARTICIPANT_DEMOTED, {
                    id: body,
                  });
                }
              }
              break;
            case GroupVideo.MESSAGE_RECEIVED:
              this.dispatchEvent(GroupVideo.MESSAGE_RECEIVED, body);
              break;
            case GroupVideo.WE_ARE_DEMOTED:
              if (this._whoAmIMeta) {
                this._whoAmIMeta.isParticipating = false;
              }
              this.dispatchEvent(GroupVideo.WE_ARE_DEMOTED);
              break;
            case GroupVideo.WE_ARE_PROMOTED:
              if (this._whoAmIMeta) {
                this._whoAmIMeta.isParticipating = true;
              }
              this.dispatchEvent(GroupVideo.WE_ARE_PROMOTED);
              break;
            default:
              break;
          }
        } catch (e) {}
      };
  
      this.iframe = iframe;
  
      window.addEventListener("message", this._onMessage);
    }
  
    /**
     * Given a participant ID (`id`), promote the participant to display its
     * video.
     * @param {string} id The ID of the participant.
     */
    promoteParticipant(id) {
      this.iframe.contentWindow.postMessage(
        JSON.stringify({
          type: "PROMOTE_PARTICIPANT",
          data: { id },
        }),
        new URL(this.iframe.src).origin
      );
    }
  
    /**
     * Given a participant ID (`id`), demote the participant away from the call,
     * and just have that participant be a passive audience.
     * @param {string} id The ID of the participant to demote.
     */
    demoteParticipant(id) {
      this.iframe.contentWindow.postMessage(
        JSON.stringify({
          type: "DEMOTE_PARTICIPANT",
          data: { id },
        }),
        new URL(this.iframe.src).origin
      );
    }
  
    /**
     * Sends a message to mute all participants, if admin
     */
    muteAllParticipants() {
      this.iframe.contentWindow.postMessage(
        JSON.stringify({
          type: "MUTE_ALL",
        }),
        new URL(this.iframe.src).origin
      );
    }
  
    /**
     * Sends a message to the intended participant, as indicated by the
     * participant ID.
     * @param {string} to The participant ID to send the message to
     * @param {any} message The message body. Could be just about any valid JSON
     *   Object.
     */
    sendMessage(to, message) {
      this.iframe.contentWindow.postMessage(
        JSON.stringify({
          type: "SEND_MESSAGE",
          data: {
            to,
            data: message,
          },
        }),
        new URL(this.iframe.src).origin
      );
    }
  
    /**
     * Note: the returning "map" object is readonly!
     */
    get participants() {
      return new Map(this._participants);
    }
  
    /**
     * A promise that represents the participant's ID.
     */
    getMyId() {
      if (this._participantId) {
        return this._participantId;
      }
  
      this._participantId = new Promise((resolve) => {
        let responded = false;
        let commandId = nextCount();
  
        const sendMessage = () => {
          setTimeout(() => {
            if (responded) {
              return;
            }
            try {
              this.iframe.contentWindow.postMessage(
                JSON.stringify({
                  type: "GET_CLIENT_ID",
                  data: { commandId },
                }),
                new URL(this.iframe.src).origin
              );
            } catch (e) {}
  
            sendMessage();
          }, 16);
        };
  
        const onMessage = (e) => {
          if (responded) {
            return;
          }
          try {
            const data = JSON.parse(e.data);
            if (data.type === "GOT_CLIENT_ID") {
              if (data.data.commandId === commandId) {
                responded = true;
                resolve(data.data.clientId);
                window.removeEventListener("message", onMessage);
              }
            }
          } catch (e) {}
        };
  
        window.addEventListener("message", onMessage);
        sendMessage();
      });
  
      return this._participantId;
    }
  
    /**
     * A promise that represent's all the user's current metadata.
     */
    getWhoAmI() {
      if (this._whoAmI) {
        return this._whoAmI;
      }
      this._whoAmI = sendRequest(
        this.iframe,
        "WHO_AM_I",
        {},
        "GOT_WHO_AM_I"
      ).then((data) => {
        this._whoAmIMeta = data;
        return this._whoAmIMeta;
      });
      return this._whoAmI;
    }
  
    /**
     * Just some cleanup (for now, just remove the `message` event listener).
     */
    destroy() {
      window.removeEventListener("message", this._onMessage);
    }
  
    /**
     * In some web browsers an iframe that requests permission to use the web
     * camera, and microphone will NOT actually prompt the user for those
     * permissions. The parent page (which the iframe is is in) MUST request those
     * permissions on behalf of the iframe, and the iframe will inherent those
     * specific permissions if the iframe has an allow="camera;microphone"
     * attribute on it.
     *
     * This function is a helper function for the parent page, to request those
     * permissions.
     */
    static async requestPermissions() {
      try {
        await window.navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        });
        return true;
      } catch (e) {
        return false;
      }
    }
  };