"use client";
import { useEffect } from "react";
import { useUserContext } from "@/context/UserContext";
import { useRouter } from "next/navigation";
import { HTTP, WS } from "@/ip";

export default function PlayerPage() {
  const userContext = useUserContext();
  const { token, name, buttonReady } = userContext.state;
  const router = useRouter();
  const ws = new WebSocket(WS("buzz"));

  useEffect(() => {
    if (!token || !name) {
      router.push("/");
    }
    const body = JSON.stringify({ name, token });
    fetch(HTTP("verify"), {
      method: "POST",
      body: body,
    })
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        if (!data.success) {
          router.push("/");
        }
        userContext.set({ buttonReady: data.buttonReady == "false" });
      })
      .catch((err) => console.error(err));
  }, [token, name]);

  // create websocket
  useEffect(() => {
    ws.onopen = () => {
      console.log("connected");
      ws.send(token);
    };
    ws.onmessage = (e) => {
      console.log(e.data);
      userContext.set({
        buttonReady: JSON.parse(e.data).buttonReady == "false",
      });
      console.log(buttonReady);
    };
    ws.onclose = () => {
      console.log("disconnected");
      router.push("/");
    };
  }, [ws]);

  return (
    <main className="flex flex-col items-center justify-center h-screen bg-blue-200">
      <h1 className="p-2 m-4 text-4xl">
        Welcome{" "}
        <span className="underline" id="name">
          {name}
        </span>
      </h1>
      <button
        onClick={() => {
          try {
            ws.send(name);
            userContext.set({ buttonReady: false });
          } catch (e) {
            console.error(e);
          }
        }}
        disabled={!buttonReady}
        style={{
          backgroundColor: buttonReady ? "lightgreen" : "gray",
          color: buttonReady ? "black" : "darkgray",
        }}
        className="p-6 w-56 h-56 select-none cursor-pointer flex flex-col justify-center items-center text-center text-6xl rounded-full transition-all duration-75"
      >
        Buzz
      </button>
    </main>
  );
}
