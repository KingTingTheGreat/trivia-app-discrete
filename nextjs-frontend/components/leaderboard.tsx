"use client";
import { useState, useEffect } from "react";
import { WS } from "@/ip";

const Leaderboard = () => {
  const [data, setData] = useState([]);
  const [ws, setWs] = useState(new WebSocket(WS("leaderboard")));

  useEffect(() => {
    ws.onopen = () => {
      console.log("connected to leaderboard");
    };
    ws.onmessage = (e) => {
      console.log("leaderboard");
      console.log(e.data);
      setData(JSON.parse(e.data));
    };
    ws.onclose = () => {
      console.log("disconnected from leaderboard");
      setData([])
      // try to reconnect
      setTimeout(() => {
        setWs(new WebSocket(WS("leaderboard")));
      }, 100);
    };
  }, [ws]);

  return (
    <div className="flex flex-col items-center">
      <h2>Leaderboard</h2>
      <table className="border-collapse">
        <thead>
          <tr>
            <th className="p-2 border-solid border-2">Rank</th>
            <th className="p-2 border-solid border-2">Name</th>
            <th className="p-2 border-solid border-2">Score</th>
          </tr>
        </thead>
        <tbody>
          {data &&
            data.map((player, index) => {
              const rank = index + 1;
              const color =
                rank === 1
                  ? "#d4af37"
                  : rank === 2
                  ? "#c0c0c0"
                  : rank === 3
                  ? "#cd7f32"
                  : "";
              return (
                <tr key={index}>
                  <td className="text-center p-2 border-2" style={{ color }}>
                    {rank}
                  </td>
                  <td className="text-center p-2 border-2">{player[0]}</td>
                  <td className="text-center p-2 border-2">{player[1]}</td>
                </tr>
              );
            })}
        </tbody>
      </table>
    </div>
  );
};

export default Leaderboard;
