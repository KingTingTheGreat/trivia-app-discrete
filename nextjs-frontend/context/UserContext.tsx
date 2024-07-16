import {
  Dispatch,
  SetStateAction,
  createContext,
  useContext,
  useState,
} from "react";

const LS_KEY = "user-data";

const defaultState: UserContextState = {
  name: "",
  token: "",
  buttonReady: true,
  password: "",
};

const currentState = (): UserContextState => {
  try {
    if (typeof window === "undefined") {
      return defaultState;
    }
    const data = localStorage.getItem(LS_KEY);
    if (!data) {
      return defaultState;
    }
    return JSON.parse(data);
  } catch (e) {
    console.error(e);
    return defaultState;
  }
};

const UserContext = createContext<UserContextType | null>(null);

export const UserContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [state, setState] = useState(currentState());
  const set = (stateDiff: Partial<UserContextState>) => {
    setState((prevState) => ({ ...prevState, ...stateDiff }));
  };
  const save = () => {
    localStorage.setItem(LS_KEY, JSON.stringify(state));
  };

  return (
    <UserContext.Provider value={{ state, set, save }}>
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
  save: () => void;
};

export type UserContextState = {
  name: string;
  token: string;
  buttonReady: boolean;
  password: string;
};
