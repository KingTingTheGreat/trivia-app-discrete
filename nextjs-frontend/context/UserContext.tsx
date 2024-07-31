import {
  Dispatch,
  SetStateAction,
  createContext,
  useContext,
  useState,
} from "react";
import { LS_KEY } from "@/constants";

const defaultState: UserContextState = {
  name: "",
  token: "",
  buttonReady: true,
  password: "",
};

const currentState = (): UserContextState | null => {
  try {
    if (typeof window === "undefined") {
      return null;
    }
    const data = localStorage.getItem(LS_KEY);
    if (!data) {
      return null;
    }
    return JSON.parse(data);
  } catch (e) {
    console.error(e);
    return null;
  }
};

const UserContext = createContext<UserContextType | null>(null);

export const UserContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [state, setState] = useState(currentState() || defaultState);
  const set = (stateDiff: Partial<UserContextState>) => {
    setState((prevState) => ({ ...prevState, ...stateDiff }));
  };
  const setSave = (stateDiff: Partial<UserContextState>) => {
    const newState = { ...state, ...stateDiff };
    setState(newState);
    try {
      localStorage.setItem(LS_KEY, JSON.stringify(newState));
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <UserContext.Provider value={{ state, set, setSave }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUserContext = (): UserContextType => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error(
      "usePlayerContext must be used within a PlayerContextProvider"
    );
  }
  return context;
};

export type UserContextType = {
  state: UserContextState;
  set: (stateDiff: Partial<UserContextState>) => void;
  setSave: (stateDiff: Partial<UserContextState>) => void;
};

export type UserContextState = {
  name: string;
  token: string;
  buttonReady: boolean;
  password: string;
};
