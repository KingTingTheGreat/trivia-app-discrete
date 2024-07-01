import { Dispatch, SetStateAction, createContext, useContext, useState } from "react";

const defaultState: UserContextState = {
	name: "",
	token: "",
	setName: () => {},
	setToken: () => {},
	buttonReady: true,
	password: "",
};

const UserContext = createContext<UserContextType | null>(null);

export const UserContextProvider = ({ children }: { children: React.ReactNode }) => {
	const [state, setState] = useState(defaultState);

	return <UserContext.Provider value={{ state, setState }}>{children}</UserContext.Provider>;
};

export const useUserContext = (): UserContextType => {
	const context = useContext(UserContext);
	if (!context) {
		throw new Error("usePlayerContext must be used within a PlayerContextProvider");
	}
	return context;
};

export type UserContextType = {
	state: UserContextState;
	setState: Dispatch<SetStateAction<UserContextState>>;
};

export type UserContextState = {
	name: string;
	token: string;
	setName: Dispatch<SetStateAction<string>>;
	setToken: Dispatch<SetStateAction<string>>;
	buttonReady: boolean;
	password: string;
};
