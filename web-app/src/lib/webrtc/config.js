export const iceServers = [
  {
    urls: 'turn:turn1.turn.group.video:3478',
    username: 'turnuser',
    credential: 'dJ4kP05PHcKN8Ubu',
  },
  {
    urls: 'turn:turn2.turn.group.video:3478',
    username: 'turnuser',
    credential: 'XzfVP8cpNEy17hws',
  },
  {
    urls: 'turns:turn1.turn.group.video:443',
    username: 'turnuser',
    credential: 'dJ4kP05PHcKN8Ubu',
  },
  {
    urls: 'turns:turn2.turn.group.video:443',
    username: 'turnuser',
    credential: 'XzfVP8cpNEy17hws',
  },
  {
    url: "stun:65.21.156.193:3478",
    username: "testyu",
    credential: "testyu123",
  },
  {
    url: "turn:65.21.156.193:3478",
    username: "testyu",
    credential: "testyu123",
  },

  //METERED TRUN

  {
    urls: "stun:stun.relay.metered.ca:80",
  },
  {
    urls: "turn:global.relay.metered.ca:80",
    username: "7d59d37ae8ff49aff498bf5a",
    credential: "+09PaL14gBF2AGcN",
  },
  {
    urls: "turn:global.relay.metered.ca:80?transport=tcp",
    username: "7d59d37ae8ff49aff498bf5a",
    credential: "+09PaL14gBF2AGcN",
  },
  {
    urls: "turn:global.relay.metered.ca:443",
    username: "7d59d37ae8ff49aff498bf5a",
    credential: "+09PaL14gBF2AGcN",
  },
  {
    urls: "turns:global.relay.metered.ca:443?transport=tcp",
    username: "7d59d37ae8ff49aff498bf5a",
    credential: "+09PaL14gBF2AGcN",
  },
]
