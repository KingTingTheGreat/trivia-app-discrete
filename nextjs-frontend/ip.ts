const IP: string = "localhost";

export const HOME = (): string => `http://${IP}:3000`;
export const HTTP = (route: string): string => `http://${IP}:8080/${route}`;
export const WS = (route: string): string => `ws://${IP}:8080/${route}`;
