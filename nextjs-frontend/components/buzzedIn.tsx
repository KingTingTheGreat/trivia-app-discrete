"use client";
import { useState, useEffect } from "react";
import { WS } from "@/ip";

const BuzzedIn = () => {
  const [data, setData] = useState([]);
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

  return (
    <div className="flex flex-col items-center">
      <h2>Buzzed In</h2>
      <table className="border-collapse">
        <thead>
          <tr>
            <th className="p-2 border-solid border-2">Name</th>
            <th className="p-2 border-solid border-2">Time</th>
          </tr>
        </thead>
        <tbody>
          {data &&
            data.map((player) => {
              return (
                <tr className="border-2" key={player}>
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

export default BuzzedIn;
