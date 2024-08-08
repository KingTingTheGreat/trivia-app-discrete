const IP: string = "host-ip-address";

export const HOME = (ip: string): string => `http://${ip}:3000`;
export const HTTP = (ip: string, route: string): string =>
  `http://${ip}:8080/${route}`;
export const WS = (ip: string, route: string): string =>
  `ws://${ip}:8080/${route}`;
