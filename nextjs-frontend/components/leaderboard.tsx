"use client";
import { useState, useEffect } from "react";
import { WS } from "@/ip";
import GameContent, { TableRow, TableData } from "./gameContent";
import { Player } from "@/types";

const Leaderboard = () => {
  const [data, setData] = useState<Player[]>([]);
  const [ws, setWs] = useState(new WebSocket(WS("leaderboard")));

  useEffect(() => {
    ws.onopen = () => {
      console.log("connected to leaderboard");
    };
    ws.onmessage = (e) => {
      console.log("leaderboard", e.data);
      setData(JSON.parse(e.data));
    };
    ws.onclose = () => {
      console.log("disconnected from leaderboard");
      setData([]);
      // try to reconnect
      setTimeout(() => {
        setWs(new WebSocket(WS("leaderboard")));
      }, 100);
    };
  }, [ws]);

  const leaderboardMapFunc = (player: Player, index: number): React.ReactNode => {
    const rank = index + 1;
    const color = rank === 1 ? "#d4af37" : rank === 2 ? "#c0c0c0" : rank === 3 ? "#cd7f32" : "";
    return (
      <TableRow index={index}>
        <TableData style={{ color }}>{rank}</TableData>
        <TableData>{player.Name}</TableData>
        <TableData>{player.Score}</TableData>
      </TableRow>
    );
  };

  return (
    <GameContent
      title="Leaderboard"
      headers={["Rank", "Name", "Score"]}
      content={data}
      mapFunc={leaderboardMapFunc}
    />
  );
};

export default Leaderboard;
