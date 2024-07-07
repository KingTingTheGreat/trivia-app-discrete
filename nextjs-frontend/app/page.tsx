"use client";
import { useRouter } from "next/navigation";
import { UserContextType, useUserContext } from "@/context/UserContext";
import { useState } from "react";
import { HTTP } from "@/ip";
import { DEFAULT_ERROR } from "@/constants";

export default function HomePage() {
  const userContext = useUserContext();
  const { name, token } = userContext.state;
  const router = useRouter();
  const [errorMessage, setErrorMessage] = useState<string>(DEFAULT_ERROR);

  const handleSubmit = async () => {
    fetch(HTTP("token"), {
      method: "POST",
      body: JSON.stringify({
        name: name.trim(),
        token,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        if (data.success == "true") {
          userContext.setState((prev) => ({
            ...prev,
            name: data.name,
            token: data.token,
          }));
          router.push("/player");
        } else {
          setErrorMessage(data.message);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <main className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-2xl font-semibold m-4">Enter your name: {name}</h1>
      <form
        className="w-fit flex justify-center align-center m-2"
        onSubmit={async (e) => {
          console.log("submitting");
          e.preventDefault();
          await handleSubmit();
        }}
      >
        <input
          required
          placeholder="Name"
          type="text"
          className="p-4 m-2 border-2 border-black rounded-lg"
          onChange={(e) => {
            userContext.setState((prev) => ({ ...prev, name: e.target.value }));
            setErrorMessage(DEFAULT_ERROR);
          }}
          value={name}
        />
        <button
          className="p-4 m-2 text-4xl cursor-pointer rounded-lg"
          style={{ backgroundColor: name === "" ? "gray" : "lightgreen" }}
          type="submit"
          disabled={name === ""}
        >
          →
        </button>
      </form>
      <p
        style={{
          visibility: errorMessage !== DEFAULT_ERROR ? "visible" : "hidden",
          color: "red",
        }}
      >
        {errorMessage}
      </p>
    </main>
  );
}
