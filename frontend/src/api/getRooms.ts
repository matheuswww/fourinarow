import { basePath } from "./path";

export async function getRooms(): Promise<any[]> {
  try {
    const response = await fetch(basePath+"/rooms", {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Error fetching rooms: ${response.statusText}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching rooms:', error);
    throw error;
  }
}
