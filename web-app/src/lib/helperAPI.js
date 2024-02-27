import PocketBase from "pocketbase";

// Define data classes for hosts, css, and rooms
export class HostData {
  constructor(name, id) {
    this.name = name;
    this.id = id;
  }
}

export class CSSData {
  constructor(id, name, style, hostId) {
    this.id = id;
    this.name = name;
    this.style = style;
    this.hostId = hostId;
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

  createRoom = async (roomData) => {
    const formattedRoomData = {
      id: roomData.id,
      name: roomData.name,
      description: roomData.description,
      hostId: roomData.hostId,
      thumbnail: roomData.thumbnail,
    };

    try {
      const response = await this.pocketBase
        .collection("rooms")
        .create(formattedRoomData);
      return response;
    } catch (error) {
      return error.data;
    }
  };
}
