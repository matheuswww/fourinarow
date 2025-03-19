import { useState, useEffect } from 'react';
import styles from './getRooms.module.css';
import { getRooms } from '@/api/getRooms';
import { connectRoom } from '@/api/connectRoom';
import { createRoom } from '@/api/createRoom';

interface RoomComponentProps {
  setIsConnected: (connected: boolean) => void;
  setOwner: (owner: boolean) => void;
  setWs: (ws: WebSocket | null) => void;
}

export default function GetRooms({ setIsConnected, setWs, setOwner }: RoomComponentProps) {
  const [rooms, setRooms] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const handleCreateRoom = async () => {
    try {
      setLoading(true);
      const ws = await createRoom();
      setWs(ws); 
      setIsConnected(true);
      setOwner(true)
    } catch (error) {
      console.error("Error creating room:", error);
    } finally {
      setLoading(false);
    }
  };

  const loadRooms = async () => {
    try {
      const roomsData = await getRooms();
      setRooms(roomsData);
    } catch (error) {
      console.error("Error loading rooms:", error);
    }
  };

  const handleConnectRoom = async (roomId: string) => {
    try {
      const ws = await connectRoom(roomId);
      setWs(ws); 
      setIsConnected(true);
    } catch (error) {
      console.error("Error connecting to room:", error);
    }
  };

  useEffect(() => {
    loadRooms();
  }, []);

  return (
    <div className={styles.container}>
      <div className={styles.createRoom}>
        <button onClick={handleCreateRoom} disabled={loading} className={styles.button}>
          {loading ? "Criando" : "Criar sala"}
        </button>
      </div>
      <div className={styles.roomList}>
        <h2>Salas disponiveis</h2>
        <ul>
          {rooms?.map((room: any) => (
            <li key={room.id}>
              <span>{room.name}</span>
              <button onClick={() => handleConnectRoom(room.id)} className={styles.connectButton}>
                Connect
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
