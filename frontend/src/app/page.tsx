"use client";
import GetRooms from "@/components/room/getRooms";
import Room from "@/components/room/room";
import { useState } from "react";

export default function Home() {
  const [owner, setOwner] = useState<boolean>(false)
  const [isConnected, setIsConnected] = useState<boolean>(false)
  const [ws, setWs] = useState<WebSocket | null>(null)
  return (
    isConnected ? <Room ws={ws} owner={owner} /> : <GetRooms setOwner={setOwner} setWs={setWs} setIsConnected={setIsConnected} />
  );
}
