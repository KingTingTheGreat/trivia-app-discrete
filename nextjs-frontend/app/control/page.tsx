"use client";
import { useState, useEffect } from "react";
import Leaderboard from "@/components/leaderboard";
import { HTTP, WS } from "@/ip";
import { defaultError } from "../page";

export default function ControlPage() {
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState(defaultError);
  const [qnum, setQnum] = useState(0);
  const [players, setPlayers] = useState<string[]>([]);
  const [playerTokens, setPlayerTokens] = useState<{ [name: string]: string }>(
    {}
  );
  const [amount, setAmount] = useState("0");
  const [name, setName] = useState("");
  const [token, setToken] = useState("");
  const [ws, setWs] = useState(new WebSocket(WS("players")));

  // create websocket
  useEffect(() => {
    ws.onopen = () => {
      console.log("connected to players");
    };
    ws.onmessage = (e) => {
      console.log("players", e.data);
      const playersList = JSON.parse(e.data) || [];
      setPlayers(playersList);

      const playerTokens: { [name: string]: string } = {};
      for (const playerData of playersList) {
        playerTokens[playerData[0]] = playerData[1];
      }
      setPlayerTokens(playerTokens);
    };
    ws.onclose = () => {
      console.log("disconnected from players");
      setPlayers([]);
      setAmount("0");
      setName("");
      setToken("");
      // try to reconnect
      setTimeout(() => {
        setWs(new WebSocket(WS("players")));
      }, 100);
    };
  }, [ws]);

  return (
    <main className="flex flex-col items-center p-2">
      <p>{name}</p>
      <p>{token}</p>
      <h1>Control Panel</h1>
      <div>
        <h4>Password</h4>
        <input
          className="m-1 p-1 border-2 border-black"
          type="password"
          placeholder="Password"
          onChange={(e) => {
            setErrorMessage(defaultError);
            setPassword(e.target.value);
          }}
        />
        <p
          className="text-md text-center"
          style={{
            visibility: errorMessage !== defaultError ? "visible" : "hidden",
            color: "red",
          }}
        >
          {errorMessage}
        </p>
      </div>
      <h4 className="p-2 m-1">
        Question #<span id="qnum">1</span>
      </h4>
      <div className="flex flex-col sm:flex-row">
        <div className="flex flex-col m-2 p-3">
          <h4>Update User Score</h4>
          <select
            className="border-2 border-black p-1"
            value={name}
            onChange={(e) => {
              console.log(playerTokens);
              const name = e.target.value;
              setName(name);
              setToken(playerTokens[name]);
              setAmount("0");
              setErrorMessage(defaultError);
            }}
          >
            <option></option>
            {players &&
              players.map((player) => {
                return (
                  <option key={player[0]} value={player[0]}>
                    {player[0]}
                  </option>
                );
              })}
          </select>
          <input
            className="border-2 border-black p-1"
            type="number"
            placeholder="Amount"
            value={amount}
            id="amount"
            onChange={(e) => setAmount(e.target.value)}
          />
          <button
            onClick={() => {
              console.log(name, token);
              fetch(HTTP("player"), {
                method: "PUT",
                body: JSON.stringify({
                  password,
                  amount,
                  token,
                  name,
                }),
              })
                .then((res) => res.json())
                .then((data) => {
                  if (data.success == "false") {
                    setErrorMessage(data.message);
                    setName("");
                  }
                })
                .catch((err) => console.error(err));
            }}
          >
            Update Score
          </button>
        </div>
        <div className="flex flex-col m-2 p-3">
          <h4>Question Controls</h4>
          <div className="flex w-full">
            <button
              className="cursor-pointer flex-1 m-2 p-2 w-fit bg-rose-200"
              onClick={() => {
                console.log("prev");
              }}
            >
              Prev
            </button>
            <button
              className="cursor-pointer flex-1 m-2 p-2 w-fit bg-green-300"
              onClick={() => {
                console.log("next");
              }}
            >
              Next
            </button>
          </div>
          <button
            className="cursor-pointer m-2 p-2 w-fit bg-stone-300 flex-1"
            onClick={() =>
              fetch(HTTP("reset"), {
                method: "POST",
                body: JSON.stringify({ password }),
              })
                .then((res) => res.json())
                .then((data) => {
                  if (data.success == "false") {
                    setErrorMessage(data.message);
                  }
                })
                .catch((err) => console.error(err))
            }
          >
            Reset Buzzers
          </button>
        </div>
      </div>
      <Leaderboard />
      <div className="flex flex-col">
        <h4>Remove Player</h4>
        <select
          className="border-2 border-black p-1"
          name="players"
          id="remove-players"
        >
          {players &&
            players.map((player) => (
              <option key={player[0]} value={player[0]}>
                {player[0]}
              </option>
            ))}
        </select>
        <button
          className="border-2 border-black p-1"
          onClick={() =>
            fetch(HTTP("player"), {
              method: "DELETE",
              body: JSON.stringify({ password }),
            })
              .then((res) => res.json())
              .then((data) => {
                if (data.success == "false") {
                  setErrorMessage(data.message);
                  setName("");
                }
              })
              .catch((err) => console.error(err))
          }
        >
          Remove Player
        </button>
      </div>
    </main>
  );
}
