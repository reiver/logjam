const myPeerConnectionConfig = {
    iceServers: [
        // {'urls': 'stun:stun.l.google.com:19302'},
        // {"urls": "turn:numb.viagenie.ca", "username":"m.mirsamie@gmail.com", "credential":"159951"},
        // {
        //   url: 'stun:stun.1und1.de:3478'
        // },
        // {
        //   url: 'stun:stun.gmx.net:3478'
        // },
        {
            url: 'stun:stun.l.google.com:19302'
        },
        {
            url: 'stun:stun1.l.google.com:19302'
        },
        {

            url: 'stun:stun2.l.google.com:19302'
        },
        {

            url: 'stun:stun3.l.google.com:19302'
        },
        {

            url: 'stun:stun4.l.google.com:19302'
        },
        /*
        {
          url: "turn:numb.viagenie.ca",
          credential: "muazkh",
          username: "webrtc@live.com",
        },
        {
          url: "turn:192.158.29.39:3478?transport=udp",
          credential: "JZEOEt2V3Qb0y27GRntt2u2PAYA=",
          username: "28224511:1379330808",
        },
        {
          url: "turn:192.158.29.39:3478?transport=tcp",
          credential: "JZEOEt2V3Qb0y27GRntt2u2PAYA=",
          username: "28224511:1379330808",
        },
        {
          url: "turn:turn.bistri.com:80",
          credential: "homeo",
          username: "homeo",
        },
        {
          url: "turn:turn.anyfirewall.com:443?transport=tcp",
          credential: "webrtc",
          username: "webrtc",
        },
        {
          url:"turn:13.250.13.83:3478?transport=udp",
          username: "YzYNCouZM1mhqhmseWk6",
          credential: "YzYNCouZM1mhqhmseWk6"
        }
        */
        /*
        {
          url: "turn:45.149.77.155:3478",
          username: "admin",
          credential: "mmcomp",
        },
        */
        {
            url: "turn:turn1.turn.group.video:3478",
            username: "turnuser",
            credential: "dJ4kP05PHcKN8Ubu",
        },
        {
            url: "turn:turn2.turn.group.video:3478",
            username: "turnuser",
            credential: "XzfVP8cpNEy17hws",
        },
        {
            url: "turns:turn1.turn.group.video:443",
            username: "turnuser",
            credential: "dJ4kP05PHcKN8Ubu",
        },
        {
            url: "turns:turn2.turn.group.video:443",
            username: "turnuser",
            credential: "XzfVP8cpNEy17hws",
        },
    ],
};