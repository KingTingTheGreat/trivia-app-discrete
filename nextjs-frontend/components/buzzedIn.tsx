"use client";
import { useState, useEffect } from "react";
import { WS } from "@/ip";
import { useUserContext } from "@/context/UserContext";

const BuzzedIn = () => {
  const userContext = useUserContext();
  const [data, setData] = useState([]);
  const [ws, setWs] = useState(
    new WebSocket(WS(userContext.state.ip, "buzzed-in"))
  );
  const buzzer = new Audio("/buzzer.mp3");

  useEffect(() => {
    ws.onopen = () => {
      console.log("connected");
    };
    ws.onmessage = (e) => {
      console.log(e.data);
      setData(JSON.parse(e.data));
    };
    ws.onclose = () => {
      console.log("disconnected");
      // try to reconnect
      setTimeout(() => {
        setWs(new WebSocket(WS(userContext.state.ip, "buzzed-in")));
      }, 1000);
    };
  }, [ws]);

  useEffect(() => {
    if (data.length > 0) {
      buzzer.play();
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
            data.map((player, index) => {
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
