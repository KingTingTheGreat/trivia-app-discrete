"use client";
import { useState, useEffect } from "react";
import { WS } from "@/ip";
import GameContent, { TableRow, TableData } from "./gameContent";
import { Player } from "@/types";

const BuzzedIn = () => {
  const [data, setData] = useState<Player[]>([]);
  const [ws, setWs] = useState(new WebSocket(WS("buzzed-in")));
  const [buzzer, setBuzzer] = useState<HTMLAudioElement | null>(null);

  useEffect(() => {
    setBuzzer(new Audio("/buzzer.mp3"));
  }, []);

  useEffect(() => {
    ws.onopen = () => {
      console.log("connected to buzzed in");
    };
    ws.onmessage = (e) => {
      console.log("buzzed in", e.data);
      setData(JSON.parse(e.data));
    };
    ws.onclose = () => {
      console.log("disconnected from buzzed in");
      setData([]);
      // try to reconnect
      setTimeout(() => {
        setWs(new WebSocket(WS("buzzed-in")));
      }, 100);
    };
  }, [ws]);

  useEffect(() => {
    if (data.length > 0) {
      buzzer?.play();
    }
  }, [data]);

  const buzzedInMapFunc = (player: Player, index: number): React.ReactNode => {
    return (
      <TableRow index={index}>
        <TableData>{player.Name}</TableData>
        <TableData>{player.Time}</TableData>
      </TableRow>
    );
  };

  return (
    <GameContent
      title="Buzzed In"
      headers={["Name", "Time"]}
      content={data}
      mapFunc={buzzedInMapFunc}
    />
  );
};

export default BuzzedIn;
