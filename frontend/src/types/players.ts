export type PlayerType = "offline" | "online";

export type Player<TStatus extends PlayerType = "online"> =
  TStatus extends "online" ? OnlinePlayer : OfflinePlayer;

export type OfflinePlayer = {
  id: number;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string | null;
  playerId: string;
  playerSteamId: string;
  name: string;
  ip: string;
  pbGuid: string;
};

export type OnlinePlayerDump = Pick<
  OfflinePlayer,
  "ip" | "pbGuid" | "name" | "playerId" | "playerSteamId"
> & {
  xVer: string;
  qPort: number;
  challenge: number;
  protocol: number;
  cgPredictItems: number;
  clAnonymous: boolean;
  clPunkbuster: boolean;
  clVoice: boolean;
  clWwwDownload: boolean;
  rate: number;
  snaps: number;
};

export type OnlinePlayer = Pick<OnlinePlayerDump, "rate" | "qPort"> & {
  id: number;
  score: number;
  ping: number;
  uuid: string;
  steamId: string;
  name: string;
  strippedName: string;
  address: string;
};
