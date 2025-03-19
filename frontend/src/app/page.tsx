"use client";
import GetRooms from "@/components/room/getRooms";
import Room from "@/components/room/room";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function Home() {
  const router = useRouter()
  const [owner, setOwner] = useState<boolean>(false)
  const [isConnected, setIsConnected] = useState<boolean>(false)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [load, setLoad] = useState<boolean>(false)

  useEffect(() => {
    if(window.localStorage.getItem("token") && window.localStorage.getItem("user_id")) {
      setLoad(true)
      return
    }
    router.push("/signup")
  }, [])

  return (
    load && (isConnected ? <Room ws={ws} owner={owner} /> : <GetRooms setOwner={setOwner} setWs={setWs} setIsConnected={setIsConnected} />)
  )
}
