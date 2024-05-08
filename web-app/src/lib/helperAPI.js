import { css } from "@emotion/react";
import PocketBase from "pocketbase";

// Define data classes for hosts, css, and rooms
export class HostData {
  constructor(name, id) {
    this.name = name;
    this.id = id;
  }
}

export class CSSData {
  constructor(id, name, style, hostId, fileHash) {
    this.id = id;
    this.name = name;
    this.style = style;
    this.hostId = hostId;
    this.fileHash = fileHash;
  }
}

export class RoomData {
  constructor(name, description, thumbnail, hostId, id) {
    this.name = name;
    this.description = description;
    this.thumbnail = thumbnail;
    this.hostId = hostId;
    this.id = id;
  }
}

export class walletData {
  constructor(address, message, signature) {
    this.address = address;
    this.message = message;
    this.signature = signature;
  }
}

// Function to convert RoomData instance to FormData
export function convertRoomDataToFormData(roomData) {
  const formData = new FormData();

  // Assuming the properties of RoomData are directly accessible
  formData.append("name", roomData.name);
  formData.append("description", roomData.description);
  formData.append("thumbnail", roomData.thumbnail); // Note: This should be a File object if you're uploading a file
  formData.append("hostId", roomData.hostId);
  formData.append("id", roomData.id);

  return formData;
}

// PocketBase manager

export class PocketBaseManager {
  constructor() {
    // Initialize PocketBase with your base URL
    this.pocketBase = new PocketBase("https://pb.greatape.stream");
  }

  // Function to create a host entry
  createHost = async (hostData) => {
    const formattedHostData = {
      name: hostData.name,
      id: hostData.id,
    };
    try {
      console.log("HOST DATA: ", formattedHostData);
      const response = await this.pocketBase
        .collection("hosts")
        .create(formattedHostData);
      return response;
    } catch (error) {
      return error.data;
    }
  };

  createCSS = async (cssData) => {
    const formattedCSSData = {
      style: cssData.style,
      name: cssData.name,
      hostId: cssData.hostId,
      id: cssData.id,
      fileHash: cssData.fileHash,
      lastUsed: false,
    };

    try {
      const response = await this.pocketBase
        .collection("css")
        .create(formattedCSSData);
      return response;
    } catch (error) {
      return error.data;
    }
  };

  updateCSS = async (cssData, lastUsed) => {
    const formattedCSSData = {
      style: cssData.style,
      name: cssData.name,
      hostId: cssData.hostId,
      id: cssData.id,
      fileHash: cssData.fileHash,
      lastUsed: lastUsed,
    };

    try {
      const response = await this.pocketBase
        .collection("css")
        .update(cssData.id, formattedCSSData);
      return response;
    } catch (error) {
      return error.data;
    }
  };

  createRoom = async (roomData) => {
    // const formattedRoomData = {
    //   id: roomData.id,
    //   name: roomData.name,
    //   description: roomData.description,
    //   hostId: roomData.hostId,
    //   thumbnail: roomData.thumbnail,
    // };

    try {
      const response = await this.pocketBase
        .collection("rooms")
        .create(roomData);
      return response;
    } catch (error) {
      return error.data;
    }
  };

  //list Data

  getHostById = async (hostId) => {
    try {
      const host = await this.pocketBase
        .collection("hosts")
        .getFirstListItem(`id="${hostId}"`);
      return host;
    } catch (error) {
      return error.data;
    }
  };

  getHostByName = async (name) => {
    try {
      const host = await this.pocketBase
        .collection("hosts")
        .getFirstListItem(`name="${name}"`);
      return host;
    } catch (error) {
      return error.data;
    }
  };

  getFullListOfHosts = async () => {
    try {
      // you can also fetch all records at once via getFullList
      const records = await this.pocketBase.collection("hosts").getFullList({
        sort: "-created",
      });
      return records;
    } catch (error) {
      return error.data;
    }
  };

  getRoomBYHostID = async (hostId) => {
    try {
      const room = await this.pocketBase
        .collection("rooms")
        .getFirstListItem(`hostId="${hostId}"`);
      return room;
    } catch (error) {
      return error.data;
    }
  };

  getFullListOfRoomsBYHostId = async (hostId) => {
    try {
      const room = await this.pocketBase.collection("rooms").getFullList({
        sort: "-created",
        filter: `hostId= '${hostId}'`,
      });
      return room;
    } catch (error) {
      return error.data;
    }
  };

  getCSSbyHostId = async (hostId) => {
    try {
      const css = await this.pocketBase
        .collection("css")
        .getFirstListItem(`hostId="${hostId}"`);
      return css;
    } catch (error) {
      return error.data;
    }
  };

  getCSSbyHash = async (hash, hostId) => {
    try {
      const css = await this.pocketBase
        .collection("css")
        .getFirstListItem(`fileHash="${hash}" && hostId="${hostId}"`);
      return css;
    } catch (error) {
      return error.data;
    }
  };

  getFullListOfCssBYHostId = async (hostId) => {
    try {
      const css = await this.pocketBase.collection("css").getFullList({
        sort: "-created",
        filter: `hostId= '${hostId}'`,
      });
      return css;
    } catch (error) {
      return error.data;
    }
  };

  deleteCssRecord = async (cssId) => {
    try {
      const res = await this.pocketBase.collection("css").delete(cssId);
      return res;
    } catch (error) {
      return error.data;
    }
  };

  saveWalletAddress = async (walletData) => {
    const data = {
      address: walletData.address,
      message: walletData.message,
      signature: walletData.signature,
    };
    try {
      console.log("formattedwalletAddressData", data);
      const response = await this.pocketBase.collection("wallets").create(data);
      return response;
    } catch (error) {
      return error.data;
    }
  };
}
